$(function() {
  $('#threads').on('click', '.thread_action', handleThreadActionButtonClick);
});


// TODO: dynamically update vote count instead of reloading the page
function handleThreadActionButtonClick(e) {
  var target = $(e.target);
  target.prop('disabled', true); // stop multiple clicks
  $.ajax({
    url: '/t/' + target.val() + '/' + target.attr('endpoint'),
    type: target.attr('method'),
    success: function(result) {
      location.reload();
    }
  });
}
