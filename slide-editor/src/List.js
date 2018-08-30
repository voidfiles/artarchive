import React, { Component } from 'react';
import { ListGroupItem, FormGroup, ListGroup, InputGroup, Button, FormControl} from 'react-bootstrap';


class ListEditor extends Component {
  constructor(props) {
    super(props);
    // Don't call this.setState() here!
    this.state = { newListItem: "" };
  }
  handleChange = (e) => {
    this.setState({newListItem: e.target.value})
  }

  addItem = () => {
    var list = this.props.list || [];
    list.push(this.state.newListItem);

    this.props.onValueChange(list)
  }

  removeItem = (index) => {
    var list = this.props.list || [];
    list.splice(index, 1);
    this.props.onValueChange(list)
  }

  render = () => {
    var ls = this.props.list || [];
    var items = ls.map((item, i) => {
      return (
        <ListGroupItem key={i}>
          <button onClick={() => {this.removeItem(i)}} type="button" className="close" aria-label="Remove">
            <span aria-hidden="true">&times;</span>
          </button>
          <span> {item}</span>
        </ListGroupItem>);
    });
    return (
      <FormGroup>
        <ListGroup>
          {items}
          <ListGroupItem>
            <InputGroup>
              <FormControl
                value={this.state.newListItem}
                 onChange={this.handleChange}
                type="text"/>
              <InputGroup.Button>
                <Button onClick={this.addItem}>Add</Button>
              </InputGroup.Button>
            </InputGroup>
          </ListGroupItem>
        </ListGroup>
      </FormGroup>
    );
  };
}

export {
  ListEditor
};
