//package jp.ne.sakura.ogiekako.fm3.solver;
//
//import com.google.inject.Inject;
//import com.google.inject.assistedinject.Assisted;
//
//import jp.ne.sakura.ogiekako.fm3.base.Board;
//import jp.ne.sakura.ogiekako.fm3.base.BoardStateAnalyzer;
//import jp.ne.sakura.ogiekako.fm3.base.BoardStateAnalyzerFactory;
//import jp.ne.sakura.ogiekako.fm3.base.State;
//import jp.ne.sakura.ogiekako.fm3.rpc.proto.Rpc;
//
//import java.util.HashSet;
//import java.util.List;
//import java.util.Set;
//
//public class StandardMatedStateArbiter implements MatedStateArbiter {
//
//  State state;
//  BoardStateAnalyzerFactory boardStateAnalyzerFactory;
//
//  @Inject
//  StandardMatedStateArbiter(@Assisted State state) {
//    this.state = state;
//  }
//
//  @Override
//  public boolean isMated() {
//    Board board = state.getBoard();
//    BoardStateAnalyzer analyzer = boardStateAnalyzerFactory.create(board);
//    Set<Rpc.Cell> kingPositions = new HashSet<>();
//    List<Rpc.PiecePosition> pieceList = board.getPieceList();
//    for (Rpc.PiecePosition piecePosition : pieceList) {
//      if (piecePosition.getPiece().getType() == Rpc.Piece.Type.OU) {
//        kingPositions.add(piecePosition.)
//      }
//    }
//  }
//}
