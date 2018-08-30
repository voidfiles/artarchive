import React, { Component } from 'react';
import {Row, Col, FormGroup, ControlLabel, FormControl} from 'react-bootstrap';
import { ListEditor } from './List.js';

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

class ArtistWebsiteURL extends Component {
  handleChange = (e) => {
    this.props.onValueChange(e.target.value);
  }

  render() {
    return (
      <FormGroup>
        <ControlLabel>Website</ControlLabel>
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

function bindThatUpdate(cls, updater) {
  return (value) => {
    cls.setState((state) => {
      updater(state, value)
      return state;
    })
  }
}

class Artist extends Component {
  constructor(props) {
    super(props);
    this.state = this.props.artist;
  }
  handleChange = (updater) => {
    return (value) => {
      this.setState((state) => {
        updater(state, value)
        this.props.onValueChange(state);
        return state;
      });

    }
  }
  render() {
    return (
      <div>
        <Row>
          <Col xs={12}>
            <ArtistName name={this.state.name} onValueChange={this.handleChange((state, name) => state.name = name)}></ArtistName>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <ArtsyURL url={this.state.name.artsy_url} artist_name={this.state.name.name} onValueChange={this.handleChange((state, url) => state.artsy_url = url)}></ArtsyURL>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <ArtistWikipediaURL url={this.state.wikipedia_url} artist_name={this.state.name} onValueChange={this.handleChange((state, url) => state.wikipedia_url = url)}></ArtistWikipediaURL>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <ArtistInstagramURL url={this.state.instagram_url} onValueChange={this.handleChange((state, url) => state.instagram_url = url)}></ArtistInstagramURL>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <ArtistTwitterURL url={this.state.twitter_url} onValueChange={this.handleChange((state, url) => state.twitter_url = url)}></ArtistTwitterURL>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <ArtistWebsiteURL url={this.state.website_url} onValueChange={this.handleChange((state, url) => state.website_url = url)}></ArtistWebsiteURL>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <FormGroup>
              <ControlLabel>Feeds</ControlLabel>
              <ListEditor
                list={this.state.feeds}
                onValueChange={this.handleChange((state, feeds) => state.feeds = feeds)}
                />
            </FormGroup>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <FormGroup>
              <ControlLabel>Sites</ControlLabel>
              <ListEditor
                list={this.state.sites}
                onValueChange={this.handleChange((state, sites) => state.sites = sites)}
                />
            </FormGroup>
          </Col>
        </Row>
      </div>
    )
  }
}

export {Artist};
