var DOMReady = function(a,b,c){b=document,c='addEventListener';b[c]?b[c]('DOMContentLoaded',a):window.attachEvent('onload',a)}


DOMReady(function () {
  loadFileChunks("");

  var id_box = document.getElementById("query");
  id_box.onkeyup = function() {
    var tagString = this.value.replace(/^\W+|\W+$/, "").replace(/\W+/g, ",");
    loadFileChunks(tagString);
  }
});

function loadFileChunks(tagString) {
  var request = new XMLHttpRequest();
  request.open('GET', '/chunks?tags=' + tagString, true);
  request.onload = function() {
    if (request.status === 200) {
      riot.mount('list', { fileChunks: JSON.parse(request.responseText) });
    }
  };
  request.send();
}
