package jp.ne.sakura.ogiekako.fm3.rpc;

import jp.ne.sakura.ogiekako.fm3.rpc.proto.Rpc;

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.logging.Logger;

public class SolverServer {

  private static final Logger logger = Logger.getLogger(SolverServer.class.getName());

  public static void main(String[] args) {
    new SolverServer().run();
  }

  private void run() {
    try (ServerSocket socketServer = new ServerSocket(8080)) {
      while (true) {
        try {
          Socket socket = socketServer.accept();
          DataInputStream in = new DataInputStream(socket.getInputStream());
          DataOutputStream out = new DataOutputStream(socket.getOutputStream());
          int len = in.readInt();
          byte[] req = new byte[len];
          for (int i = 0; i < len; i++) {
            req[i] = in.readByte();
          }
          Rpc.Problem problem = Rpc.Problem.parseFrom(req);
          Rpc.SolutionSet solutionSet = solve(problem);
          byte[] res = solutionSet.toByteArray();
          logger.info(String.format("Sending solutionSet. Length = %d.", res.length));
          out.writeInt(res.length);
          out.write(res, 0, res.length);
        } catch (IOException e) {
          logger.warning(String.format("Failed to process the query. Error: %s", e));
        }
      }
    } catch (IOException e) {
      logger.warning(String.format("Failed to instantiate the server socket. Error: %s", e));
    }
  }

  private Rpc.SolutionSet solve(Rpc.Problem problem) {
    // TODO(oka): Replace with the actual implementation.
    return Rpc.SolutionSet.newBuilder().addSolution(Rpc.Solution.newBuilder()
        .addMove(Rpc.Move.newBuilder()
            .setMover(Rpc.Player.HITHER)
            .setSource(Rpc.Move.Source.getDefaultInstance())
            .setDestination(Rpc.Move.Destination.newBuilder()
                .setCell(Rpc.Cell.newBuilder().setRow(1).setCol(1))
                .setPiece(Rpc.Piece.newBuilder().setType(Rpc.Piece.Type.HI))))
        .addMove(Rpc.Move.newBuilder()
            .setMover(Rpc.Player.YONDER)
            .setSource(
                Rpc.Move.Source.newBuilder().setCell(Rpc.Cell.newBuilder().setRow(5).setCol(5)))
            .setDestination(Rpc.Move.Destination.newBuilder()
                .setCell(Rpc.Cell.newBuilder().setRow(1).setCol(1))
                .setPiece(Rpc.Piece.newBuilder().setType(Rpc.Piece.Type.KAKU))))
        .addMove(Rpc.Move.newBuilder()
            .setMover(Rpc.Player.HITHER)
            .setSource(Rpc.Move.Source.getDefaultInstance())
            .setDestination(Rpc.Move.Destination.newBuilder()
                .setCell(Rpc.Cell.newBuilder().setRow(8).setCol(8))
                .setPiece(Rpc.Piece.newBuilder().setType(Rpc.Piece.Type.KIN)))))
        .build();
  }
}
