package jp.ne.sakura.ogiekako.fm3;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.URI;
import java.util.HashMap;
import java.util.Map;

public class SimpleHttpServer {

  public static void main(String[] args) throws Exception {
    HttpServer server = HttpServer.create(new InetSocketAddress(8000), 0);
    server.createContext("/test", new MyHandler());
    server.setExecutor(null); // creates a default executor
    server.start();
  }

  static class MyHandler implements HttpHandler {

    public void handle(HttpExchange t) throws IOException {
      String response = "Welcome Real's HowTo test page";
      t.sendResponseHeaders(200, response.length());
      String method = t.getRequestMethod();
      System.err.println(method);
      queryToMap(t.getRequestURI().getQuery());
      URI requestURI = t.getRequestURI();
      System.err.println(requestURI.getQuery());
      OutputStream os = t.getResponseBody();
      os.write(response.getBytes());
      os.close();
    }

    public Map<String, String> queryToMap(String query) {
      Map<String, String> result = new HashMap<>();
      for (String param : query.split("&")) {
        String pair[] = param.split("=");
        if (pair.length > 1) {
          result.put(pair[0], pair[1]);
        } else {
          result.put(pair[0], "");
        }
      }
      return result;
    }
  }
}
