$(document).ready(function() {
  $('#threads').on('click', '.thread_upvote_button', function(e) {
  	handleVoteButtonClick(e, 'POST', 'thread_remove_vote_button', 'remove vote');
  });

  $('#threads').on('click', '.thread_remove_vote_button', function(e) {
  	handleVoteButtonClick(e, 'DELETE', 'thread_upvote_button', 'upvote');
  });
});


function handleVoteButtonClick(e, call_type, new_class, new_html) {
		$(e.target).hide();
	  $.ajax({
	  	url: '/upvote/' + e.target.value,
	  	type: call_type,
	  	success: function(result) {
	      location.reload();
	      // $(e.target).attr('class', new_class)
	      // $(e.target).html(new_html);
	      // $(e.target).show();
	  	}
	  });
}
