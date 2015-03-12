//package jp.ne.sakura.ogiekako.fm3.solver;
//
//import com.google.common.collect.ImmutableList;
//import com.google.common.truth.Truth;
//import com.google.inject.Guice;
//
//import jp.ne.sakura.ogiekako.fm3.base.BaseModule;
//import jp.ne.sakura.ogiekako.fm3.base.Board;
//import jp.ne.sakura.ogiekako.fm3.base.BoardFactory;
//import jp.ne.sakura.ogiekako.fm3.rpc.proto.PieceUtil;
//import jp.ne.sakura.ogiekako.fm3.rpc.proto.Rpc;
//
//import org.junit.Before;
//import org.junit.Test;
//
//public class StandardMatedStateArbiterTest {
//
//  BoardFactory boardFactory;
//
//  @Before
//  public void setUp() throws Exception {
//    boardFactory = Guice.createInjector(new BaseModule())
//        .getInstance(BoardFactory.class);
//  }
//
//  @Test
//  public void isMated_True() throws Exception {
//    Rpc.State state = Rpc.State.newBuilder()
//        .setHeight(9)
//        .setWidth(9)
//        .setTurn(Rpc.Player.YONDER)
//        .addAllPiecePosition(ImmutableList.of(
//            PieceUtil.pieceOnBoard(5, 1, Rpc.Piece.Type.OU, Rpc.Player.YONDER),
//            PieceUtil.pieceOnBoard(5, 2, Rpc.Piece.Type.KIN, Rpc.Player.HITHER),
//            PieceUtil.pieceOnBoard(5, 3, Rpc.Piece.Type.FU, Rpc.Player.HITHER)))
//        .build();
//    Board board = boardFactory.create(state);
//    MatedStateArbiter arbiter = new StandardMatedStateArbiter(board);
//    Truth.assertThat(arbiter.isMated()).isTrue();
//  }
//
//  @Test
//  public void isMated_False() throws Exception {
//    Rpc.State state = Rpc.State.newBuilder()
//        .setHeight(9)
//        .setWidth(9)
//        .setTurn(Rpc.Player.YONDER)
//        .addAllPiecePosition(ImmutableList.of(
//            PieceUtil.pieceOnBoard(5, 1, Rpc.Piece.Type.OU, Rpc.Player.YONDER),
//            PieceUtil.pieceOnBoard(5, 2, Rpc.Piece.Type.KIN, Rpc.Player.HITHER)))
//        .build();
//    Board board = boardFactory.create(state);
//    MatedStateArbiter arbiter = new StandardMatedStateArbiter(board);
//    Truth.assertThat(arbiter.isMated()).isFalse();
//  }
//
//}