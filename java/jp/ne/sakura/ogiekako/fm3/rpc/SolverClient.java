package jp.ne.sakura.ogiekako.fm3.rpc;

import jp.ne.sakura.ogiekako.fm3.rpc.proto.PieceUtil;
import jp.ne.sakura.ogiekako.fm3.rpc.proto.Rpc;

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.util.logging.Logger;

public class SolverClient {

  private static Logger logger = Logger.getLogger(SolverClient.class.getName());

  // For test.
  public static void main(String[] args) throws IOException {
    Rpc.Problem problem = Rpc.Problem.newBuilder()
        .setState(Rpc.State.newBuilder()
            .setTurn(Rpc.Player.HITHER)
            .addPiecePosition(PieceUtil.pieceOnBoard(5, 5, Rpc.Piece.Type.KAKU, Rpc.Player.YONDER))
            .addPiecePosition(PieceUtil.pieceOnBoard(1, 9, Rpc.Piece.Type.OU, Rpc.Player.YONDER))
            .addPiecePosition(PieceUtil.pieceOnBoard(3, 7, Rpc.Piece.Type.KEI, Rpc.Player.HITHER))
            .addPiecePosition(PieceUtil.pieceOnHand(Rpc.Piece.Type.HI, Rpc.Player.HITHER))
            .addPiecePosition(PieceUtil.pieceOnHand(Rpc.Piece.Type.KIN, Rpc.Player.HITHER))
            .build())
        .setRule(Rpc.Problem.Rule.HELP)
        .setLimit(3)
        .build();

    Rpc.SolutionSet solution = new SolverClient().solve(problem);
    System.out.println(solution);
  }

  public Rpc.SolutionSet solve(Rpc.Problem problem) throws IOException {
    byte[] req = problem.toByteArray();
    try (Socket socket = new Socket("localhost", 8080)) {
      DataInputStream in = new DataInputStream(socket.getInputStream());
      DataOutputStream out = new DataOutputStream(socket.getOutputStream());
      out.writeInt(req.length);
      out.write(req, 0, req.length);
      out.flush();
      int len = in.readInt();
      logger.info(String.format("Got solution. Length = %d.", len));
      byte[] res = new byte[len];
      for (int i = 0; i < len; i++) {
        res[i] = in.readByte();
      }
      return Rpc.SolutionSet.parseFrom(res);
    } catch (IOException e) {
      e.printStackTrace();
      throw e;
    }
  }
}
