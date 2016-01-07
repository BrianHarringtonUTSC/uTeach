var TopicsContainer = React.createClass({
  loadCommentsFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    this.loadCommentsFromServer();
  },
  render: function() {
    return (
      <div className="topicsContainer">
        <h1>Topics</h1>
        <TopicList data={this.state.data} />
      </div>
    );
  }
});

var TopicList = React.createClass({
  render: function() {
    var topicNodes = this.props.data.map(function(topic) {
      return (
        <Topic name={topic.name} key={topic.id} />
      );
    });
    return (
      <div className="topicsList">
        {topicNodes}
      </div>
    );
  }
});

var Topic = React.createClass({
  render: function() {
    return (
      <div className="topic">
          {this.props.name}
      </div>
    );
  }
});

ReactDOM.render(
  <TopicsContainer url="/api/topics" />,
  document.getElementById('content')
);
