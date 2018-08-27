import React, { Component } from 'react';
import { FormGroup, ControlLabel, FormControl} from 'react-bootstrap';

class ArtistName extends Component {

  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }
  render() {
    return (
      <FormGroup>
        <ControlLabel>Name</ControlLabel>
        <FormControl
          type="text"
          value={this.props.name}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
      </FormGroup>
    );
  }
}

class ArtsyURL extends Component {
  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }

  searchUrl = () => {
    if (!this.props.artist_name) {
      return "";
    }
    let searchUrl = 'https://www.artsy.net/search?q=' + this.props.artist_name;
    return (<p><a href={searchUrl} >Search on Artsy</a></p>)
  }

  render() {
    return (
      <FormGroup>
        <ControlLabel>Artsy</ControlLabel>
        <FormControl
          type="text"
          value={this.props.url}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
        {this.searchUrl()}
      </FormGroup>
    );
  }
}

class ArtistWikipediaURL extends Component {
  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }

  searchUrl = () => {
    if (!this.props.artist_name) {
      return "";
    }
    let searchUrl = 'https://en.wikipedia.org/w/index.php?search=' + this.props.artist_name;
    return (<p><a href={searchUrl} >Search on Wikipedia</a></p>)
  }

  render() {
    return (
      <FormGroup>
        <ControlLabel>Wikipedia</ControlLabel>
        <FormControl
          type="text"
          value={this.props.url}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
        {this.searchUrl()}
      </FormGroup>
    );
  }
}

class ArtistInstagramURL extends Component {
  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }

  render() {
    return (
      <FormGroup>
        <ControlLabel>Instagram</ControlLabel>
        <FormControl
          type="text"
          value={this.props.url}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
      </FormGroup>
    );
  }
}

class ArtistTwitterURL extends Component {
  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }

  render() {
    return (
      <FormGroup>
        <ControlLabel>Twitter</ControlLabel>
        <FormControl
          type="text"
          value={this.props.url}
          placeholder="Enter text"
          onChange={this.handleChange}
        />
      </FormGroup>
    );
  }
}

export {
  ArtistName,
  ArtsyURL,
  ArtistWikipediaURL,
  ArtistInstagramURL,
  ArtistTwitterURL
};
