class App extends React.Component<{name: string}, {}> {
  render() {
    return <div>Hello, {this.props.name}</div>;
  }
}

ReactDOM.render(<App name="Willson" />, document.getElementById("root"));