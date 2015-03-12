package jp.ne.sakura.ogiekako.fm3.rpc.proto;

public class PieceUtil {

  public static Rpc.PiecePosition pieceOnBoard(int row, int col, Rpc.Piece.Type type,
      Rpc.Player owner) {
    return Rpc.PiecePosition.newBuilder()
        .setPosition(Rpc.Position.newBuilder()
            .setCell(Rpc.Cell.newBuilder().setRow(row).setCol(col))
            .setOwner(owner))
        .setPiece(Rpc.Piece.newBuilder()
            .setType(type))
        .build();
  }

  public static Rpc.PiecePosition pieceOnHand(Rpc.Piece.Type type, Rpc.Player owner) {
    return Rpc.PiecePosition.newBuilder()
        .setPosition(Rpc.Position.newBuilder()
            .setOwner(owner))
        .setPiece(Rpc.Piece.newBuilder()
            .setType(type))
        .build();
  }
}
