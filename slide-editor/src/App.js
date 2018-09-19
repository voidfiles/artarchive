import React, { Component } from 'react';
import { Navbar, Grid, Row, Col, Image, Button} from 'react-bootstrap';
import './App.css';
import Cookie from 'js-cookie';
import queryString from './query-string.js';
import {SiteTitle, SiteUrl} from './Site.js';
import {Artist} from './Artist.js';
import {Auth} from './Auth.js'

function bindThatUpdate(cls, updater) {
  return (value) => {
    cls.setState((state) => {
      updater(state, value)
      return state;
    })
  }
}

class ObjectClient extends Component {
  constructor(baseURL, token) {
    super();
    this.headers = new Headers();

    this.headers.append('Content-Type', 'application/json');
    this.headers.append('Accepts', 'application/json');
    this.headers.append('Authorization', 'Basic ' + btoa(token));
    this.baseURL = baseURL
  }

  getObject = (key) => {
    return fetch(this.baseURL + '/slides/' + key, {
     method: 'GET',
     headers: this.headers,
   }).then((response) => response.json());
  }

  saveObject = (key, data) => {
    return fetch(this.baseURL + '/slides/' + key, {
     method: 'POST',
     headers: this.headers,
     body: JSON.stringify(data)
   }).then((response) => response.json());
  }
}

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    var parsed = queryString.parse(window.location.search);
    if (!parsed.key) {
      this.state.error = "No key query param";
      return;
    }

    this.authCookie = Cookie.get('auth')
    this.objectClient = new ObjectClient("https://blooming-sands-87266.herokuapp.com", this.authCookie)
    this.key = parsed.key;
    this.objectClient.getObject(this.key).then((slide) => {
      this.setState((state) => {
        if (!slide.artists) {
          slide.artists = [{}];
        }
        state.slide = slide;
        return state;
      });
    });
  }

  handleTokenChange = (token) => {
    if (!token) {
      return;
    }
    this.authCookie = token;
    Cookie.set('auth', token);
    this.objectClient = new ObjectClient("https://blooming-sands-87266.herokuapp.com", token)
    this.forceUpdate();
  }
  saveSlide = (e) => {
    e.preventDefault();
    var slideToSave = Object.assign({}, this.state.slide);
    e.preventDefault();
    this.objectClient.updateObject(this.key, slideToSave).then((data) => {
      this.setState((state) => {
        state.slide = data;
      })
    });
  }
  render() {

    if (!this.authCookie) {
      return (
        <Auth onToken={this.handleTokenChange}></Auth>
      )
    };

    if (this.state.error) {
      return (
        <p>{this.state.error}</p>
      )
    };

    if (!this.state.slide) {
      return (
        <p>Loading...</p>
      )
    };

    var artists = this.state.slide.artists.map((artist, i) => {
        return (<Artist key={i} artist={artist} onValueChange={bindThatUpdate(this, (state, artist) => state.slide.artists[i] = artist)} />)
    });
    return (
      <div className="App">
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="#home">Slide Editor</a>
            </Navbar.Brand>
          </Navbar.Header>
        </Navbar>
        <form onSubmit={this.saveSlide}>
          <Grid>
            <Row>
              <Col xs={12}>
                <pre>
                  {JSON.stringify(this.state, null, 2)}
                </pre>
              </Col>
            </Row>

            <Row>
              <Col xs={12}>
                <h3>Site</h3>
              </Col>
            </Row>
            <Row>
              <Col xs={12} md={6}>
                <SiteTitle
                  title={this.state.slide.site.title}
                  onValueChange={bindThatUpdate(this, (state, title) => state.slide.site.title = title)}></SiteTitle>
              </Col>
              <Col xs={12} md={6}>
                <SiteUrl url={this.state.slide.site.url} onValueChange={bindThatUpdate(this, (state, url) => state.slide.site.url = url)}></SiteUrl>
              </Col>
            </Row>

            <Row>
              <Col xs={12}>
                <h3>Page</h3>
              </Col>
            </Row>
            <Row>
              <Col xs={12} md={6}>
                <p><a href={this.state.slide.page.url}>{this.state.slide.page.title}</a></p>
                <p>Found on: {new Date(this.state.slide.page.published).toString()}</p>
              </Col>
            </Row>
            <Row>
              <Col xs={12}>
                <h3>Artists</h3>
              </Col>
            </Row>
            {artists}
            <Row>
              <Col xs={12}>
                <h3>Image</h3>
              </Col>
            </Row>
            <Row>
              <Col xs={12}>
                <Image src={this.state.slide.source_image_url} responsive />
              </Col>
            </Row>
            <Button type="submit">Save</Button>
          </Grid>
        </form>
      </div>
    );
  }
}

export default App;
