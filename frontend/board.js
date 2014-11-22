/*
 * Piece names:
 * fu, kyo, kei, gin, kin, kaku, hi, ou,
 * to, nariKyo, nariKei, nariGin, uma, ryu
 */
var pieces = ["fu","kyo","kei","gin","kin","kaku","hi"]    
var yonder = {
  fu: 18,
  kyo: 4,
  kei: 4,      
  gin: 4,
  kin: 4,
  kaku: 2,
  hi: 2,
};

var hither = {fu:0,kyo:0,kei:0,gin:0,kin:0,kaku:0,hi:0};

var yonderOuOnBoard = 0;
var hitherOuOnBoard = 0;

function PieceOwner(piece, hither) {
  this.piece = piece;
  this.hither = hither
}

var board = [];
for (var i = 0; i < 9; i++) {
  row = [];
  for (var j = 0; j < 9; j++) {
    row.push(new PieceOwner("",true));
  }
  board.push(row);
}

// Position of the cursor.
var x=0,y=0;
var onBoard = true;
// Which place on hither hands is selected.
var handPlace = "";

function hasClass(elm, cls) {
  var reg = new RegExp('(\\s|^)' + cls + '(\\s|$)');
  return reg.test(elm.className);
}

function removeClass(elm, cls) {
if (hasClass(elm, cls)) {
    var reg = new RegExp('(\\s|^)' + cls + '(\\s|$)');
    elm.className = elm.className.replace(reg,' ').trim();
  }
}

function addClass(elm, cls) {
  if (!hasClass(elm, cls)) {
    elm.className = (elm.className + " " + cls).trim();
  }
}

// Move to board cell.
function moveTo(nx, ny) {
  var elm = document.getElementsByClassName("selected")[0];
  removeClass(elm, "selected");
  var newId = "C" + nx + ny;
  var newCell = document.getElementById(newId);
  addClass(newCell, "selected");
  x = nx;
  y = ny;
  onBoard = true;
}

function moveToHand(piece) {
  var elm = document.getElementsByClassName("selected")[0];
  removeClass(elm, "selected");
  var p = pieces.indexOf(piece);
  var id = "H" + p;
  elm = document.getElementById(id);
  addClass(elm, "selected");
  handPlace = piece;
  onBoard = false;
}

function moveWithKey(key) {
  console.log(x,y,onBoard,handPlace,key);
  if (onBoard) {
    var dx = 0, dy = 0;
    if (key == 37) {// left
      dy -= 1;
    } else if (key == 38) {// up
      dx -= 1;
    } else if (key == 39) {// right
      dy += 1;
    } else if (key == 40) {// down
      dx += 1;
    } else {
      return;
    }
    var nx = x + dx;
    var ny = y + dy;
    if (nx >= 9) {
      onBoard = false;
      var p = ny;
      if (p >= pieces.length) {
        p = pieces.length - 1;
      }
      moveToHand(pieces[p]);
    } else {
      if (nx < 0) nx = 0;
      if (ny < 0) ny = 0;
      if (nx >= 9) nx = 8;
      if (ny >= 9) ny = 8;
      moveTo(nx,ny);
    }
  } else {
    var p = pieces.indexOf(handPlace);
    if (key == 37) {// left
      moveToHand(pieces[Math.max(0, p-1)]); 
    } else if (key == 39) {// right
      moveToHand(pieces[Math.min(pieces.length - 1, p+1)]);
    } else if (key == 38) {// up
      moveTo(8, p);
    }
  }
}

function getPieceWithKey(key) {
  if (key == 75) { // K
    return "fu";
  }
  else if (key == 74) { // J
    return "kyo";
  }
  else if (key == 72) { // H
    return "kei";
  }
  else if (key == 71) { // G
    return "gin";
  }
  else if (key == 70) { // F
    return "kin";
  }
  else if (key == 68) { // D
    return "kaku";
  }
  else if (key == 83) { // S
    return "hi"
  }
  else if (key == 65) { // A
    return "ou";
  }
  else if (key == 73) { // I
    return "to";
  }
  else if (key == 85) { // U
    return "nariKyo";
  }
  else if (key == 89) { // Y
    return "nariKei";
  }
  else if (key == 84) { // T
    return "nariGin";
  }
  else if (key == 69) { // E
    return "uma";
  }
  else if (key == 87) { // W
    return "ryu";
  }
  else if (key == 32) {
    return "";
  }
  return null;
}

function toRaw(piece) {
  if (piece == "to") {
    return "fu";
  } else if (piece == "nariKyo") {
    return "kyo";
  } else if (piece == "nariKei") {
    return "kei";
  } else if (piece == "nariGin") {
    return "gin";
  } else if (piece == "uma") {
    return "kaku";
  } else if (piece == "ryu") {
    return "hi";
  } else {
    return piece;
  }
}

function remove(pieceOwner) {
  p = pieceOwner.piece;
  hi = pieceOwner.hither;
  if (p == "") return;
  if (p == "ou") {
    if (hi) {
      hitherOuOnBoard = false;
    } else {
      yonderOuOnBoard = false;
    }
  } else {
    p = toRaw(p);
    yonder[p]++;
  }
}

function add(pieceOwner) {
  p = pieceOwner.piece;
  hi = pieceOwner.hither;
  if (p == "") return true;
  if (p == "ou") {
    if (hi) {
      if (hitherOuOnBoard) {
        return false;
      } else {
        hitherOuOnBoard = true;
        return true;
      }
    } else {
      if (yonderOuOnBoard) {
        return false;
      } else {
        yonderOuOnBoard = true;
        return true;
      }
    }
  } else {
    p = toRaw(p);
    if (yonder[p] <= 0) {
      return false;
    } else {
      yonder[p]--;
      return true;
    }
  }
}

function putPieceOwner(cur, nxt) {
  remove(cur);
  var put = add(nxt);
  if (!put) {
    add(cur);
    return false;
  }
  return true;
}

function toHTML(piece) {
  if (piece == "fu") {
    return "歩"
  } else if (piece == "kyo") {
    return "香";
  } else if (piece == "kei") {
    return "桂";
  } else if (piece == "gin") {
    return "銀";
  } else if (piece == "kin") {
    return "金";
  } else if (piece == "kaku") {
    return "角";
  } else if (piece == "hi") {
    return "飛";
  } else if (piece == "ou") {
    return "王";
  } else if (piece == "to") {
    return "と";
  } else if (piece == "nariKyo") {
    return "杏";
  } else if (piece == "nariKei") {
    return "圭";
  } else if (piece == "nariGin") {
    return "全";
  } else if (piece == "uma") {
    return "馬";
  } else if (piece == "ryu") {
    return "竜";
  } else {
    return "";
  }
}

function putPieceToBoard(i, j, pieceOwner) {
  var cur = board[i][j];
  var nxt = pieceOwner;
  var put = putPieceOwner(cur, nxt);
  var hi = pieceOwner.hither;
  if (put) {
    board[i][j] = nxt;
    var cell = document.getElementById("C" + i + j);
    if (!hi) {
      addClass(cell, "yonder");
    } else {
      removeClass(cell, "yonder");
    }
    cell.innerHTML = ""; // Prevent buggy behaviour.
    cell.innerHTML = toHTML(pieceOwner.piece);
  }
}

function addPieceToHand(piece, delta) {
  if (yonder[piece] - delta < 0 || hither[piece] + delta < 0) {
    return;
  }
  hither[piece] += delta;
  yonder[piece] -= delta;
  elm = document.getElementById("hN" + pieces.indexOf(piece));
  elm.innerHTML = hither[piece];
}

function onKeyPress(e) {
  if (e.metaKey || e.ctrlKey || e.altKey) {
    return;
  }
  var key = e.keyCode;
  var shift = e.shiftKey;
  if (37 <= key && key <= 40) {// arrow keys
    moveWithKey(key);
  } else {
    if (onBoard) {
      if (key == 32 || 65 <= key && key <= 90) { // ' ' or 'A' - 'Z'
        var piece = getPieceWithKey(key);
        if (piece == null) {
          return;
        }
        putPieceToBoard(x, y, new PieceOwner(piece, !shift));
      }
    } else {
      console.log(key);
      var d = 0;
      if (key == 187) {// + or =
        d = 1;
      } else if (key == 189) {// - or _
        d = -1;
      } else if (key == 32) {// ' '
        d = -hither[handPlace];
      } else if (49 <= key && key <= 57) {// 1 - 9
        if (shift) {
          d = -(key - 48);
        } else {
          d = key - 48;
        }
      }
      if (d != 0) {
        addPieceToHand(handPlace, d);
      }
    }
  }
}

///// Set the initial condition. /////
window.onload = function() {
  for (var i = 0; i < 9; i++) {
    for (var j = 0; j < 9; j++) {
      (function(i, j) {
        var id = "C" + i + j;
        var cell = document.getElementById(id);
        cell.onclick = function() {
          console.log(id);
          moveTo(i,j);
        }
      })(i,j);
    }
  }
  addClass(document.getElementById("C00"), "selected");

  for (var i = 0; i < pieces.length; i++) {
    var piece = pieces[i];
    var elm = document.getElementById("H" + i);
    (function(i) {
      elm.onclick = function() {
        moveToHand(pieces[i]);
      }
    })(i);
    var name = elm.getElementsByClassName("name")[0];
    name.innerHTML = toHTML(piece);
    var num = document.getElementById("hN" + i);
    num.innerHTML = 0;
  }
  document.onkeydown = onKeyPress;

  // Parse query
  var search = decodeURIComponent(location.search);
  if (location.search.length > 0) {
    var p = search.indexOf("q=");
    var queryStr = search.substring(p + 2);
    var obj = JSON.parse(queryStr);
    document.getElementById("rule").value = obj.r;
    document.getElementById("numMove").value = obj.n;
    for (var i = 0; i < 9; i++) {
      for (var j = 0; j < 9; j++) {
        if (obj.b[i][j].piece != "") {
          putPieceToBoard(i, j, obj.b[i][j]);
        }
      }
    }
    var hi = obj.h;
    for (var i = 0; i < pieces.length; i++) {
      var p = pieces[i];
      if (hi.hasOwnProperty(p)) {
        addPieceToHand(p, hi[p]);
      }
    }
    // TODO(ogiekako): Submit the query background.
  }
}

function submit() {
  var ruleElm = document.getElementById("rule");
  var rule = ruleElm.options[ruleElm.selectedIndex].value;
  var numMove = document.getElementById("numMove").value;

  var params = {
    y: yonder,
    h: hither,
    b: board,
    r: rule,
    n: numMove
  }
  json = JSON.stringify(params);
  str = encodeURIComponent(json);
  console.log(str);
  document.location = "board.html?q=" + str;
}

