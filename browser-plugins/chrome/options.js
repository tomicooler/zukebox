function save_options() {
  var address = document.getElementById('address').value;
  var user = document.getElementById('user').value;
  chrome.storage.sync.set({
    address: address,
    user: user
  }, function() {
    var status = document.getElementById('status');
    status.textContent = 'Options saved.';
    setTimeout(function() {
      status.textContent = '';
    }, 750);
  });
}

function restore_options() {
  chrome.storage.sync.get({
    address: 'http://10.50.1.13:5000',
    user: "Stranger"
  }, function(items) {
    document.getElementById('address').value = items.address;
    document.getElementById('user').value = items.user;
  });
}
document.addEventListener('DOMContentLoaded', restore_options);
document.getElementById('save').addEventListener('click', save_options);
