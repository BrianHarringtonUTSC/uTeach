$(function() {
  $('.post-action').on('click', handlePostActionButtonClick);
});


// TODO: dynamically update vote count instead of reloading the page
function handlePostActionButtonClick(e) {
  console.log(1);
  var target = $(e.target);
  target.prop('disabled', true); // stop multiple clicks
  $.ajax({
    url: target.attr('url'),
    type: target.attr('method'),
    success: function(result) {
      location.reload();
    }
  });
}
