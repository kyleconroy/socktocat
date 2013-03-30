var connection = null;

var listenForEvents = function() {
  if (connection != null) {
    connection.close();
  }

  var owner = $("#repo input[name=owner]").val();
  var repo = $("#repo input[name=repo]").val();
  var url = 'ws://' + window.location.host + '/hooks?repo=' + owner + '/' + repo;

  $("#events").empty();

  connection = new WebSocket(url);

  // Log errors
  connection.onerror = function (error) {
    console.log('WebSocket Error ' + error);
  };

  // Log events from Github
  connection.onmessage = function (e) {
    var events = JSON.parse(e.data); // An array of Github events
    for(var i = 0; i < events.length; i++) {
      var hook = events[i];
      console.log(hook);
      $("#events").prepend(
        $("<tr>").append(
          $("<td>").text(hook.id),
          $("<td>").text(hook.type),
          $("<td>").text(hook.created_at)
        )
      )
    }
  };
}

$("#repo").change(function() {
  listenForEvents();
});

listenForEvents();
