var connection = new WebSocket('ws://socktocat.org/hooks?repo=:owner/:repo');
 
// Log events from Github
connection.onmessage = function (e) {
  var events = JSON.parse(e.data); // An array of Github events
  console.log(events);
};
