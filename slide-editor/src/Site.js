import React, { Component } from 'react';
import { FormGroup, ControlLabel, FormControl} from 'react-bootstrap';

class SiteTitle extends Component {
  constructor(props, context) {
    super(props, context);

    this.handleChange = this.handleChange.bind(this);

  }

  getValidationState() {
    return null
  }

  handleChange(e) {
    this.props.onValueChange(e.target.value);
  }
  render() {
    return (
      <FormGroup
        controlId="formBasicText"
        validationState={this.getValidationState()}
      >
        <ControlLabel>Title</ControlLabel>
        <FormControl
          type="text"
          value={this.props.title}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
        <FormControl.Feedback />
      </FormGroup>
    );
  }
}

class SiteUrl extends Component {
  constructor(props, context) {
    super(props, context);

    this.handleChange = this.handleChange.bind(this);

  }

  getValidationState() {
    return null
  }

  handleChange(e) {
    this.props.onValueChange(e.target.value);
  }
  render() {
    return (
      <FormGroup
        controlId="formBasicText"
        validationState={this.getValidationState()}
      >
        <ControlLabel>URL</ControlLabel>
        <FormControl
          type="text"
          value={this.props.url}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
        <FormControl.Feedback />
      </FormGroup>
    );
  }
}

export {SiteTitle, SiteUrl};
