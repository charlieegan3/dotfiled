var request = new XMLHttpRequest();
request.open('GET', '/filechunks', true);
request.onload = function() {
  if (request.status === 200) {
    riot.mount('list', { fileChunks: JSON.parse(request.responseText) });
  }
};
request.send();
