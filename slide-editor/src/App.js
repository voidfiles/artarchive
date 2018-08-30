import React, { Component } from 'react';
import { Navbar, Grid, Row, Col, Image, Button} from 'react-bootstrap';
import './App.css';
import Cookie from 'js-cookie';
import queryString from './query-string.js';
import {SiteTitle, SiteUrl} from './Site.js';
import {Artist} from './Artist.js';
import {Auth} from './Auth.js'
import S3 from 'aws-sdk/clients/s3';

function bindThatUpdate(cls, updater) {
  return (value) => {
    cls.setState((state) => {
      updater(state, value)
      return state;
    })
  }
}

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    var parsed = queryString.parse(window.location.search);
    if (!parsed.data) {
      this.state.error = "No data query param";
      return;
    }

    this.authCookie = Cookie.get('auth');
    if (this.authCookie) {
      this.parseAuthToken(this.authCookie);
    }

    this.s3 = new S3({
      endpoint: "http://art.rumproarious.com.s3.amazonaws.com",
      region: "us-west-2",
      accessKeyId: this.aws_access_key_id,
      secretAccessKey: this.aws_secret_access_key,
      s3BucketEndpoint: true,
    });

    this.key = parsed.data;
    this.s3.getObject({
      Bucket: "art.rumproarious.com",
      Key: this.key,
    }, (err, data) => {
      var slide = JSON.parse(data.Body);
      this.setState((state) => {
        if (!slide.artist) {
          slide.artist = {};
        }
        if (!slide.artists) {
          slide.artists = [{}];
        }
        state.slide = slide;
        return state;
      })
    });
  }

  parseAuthToken(token) {
    var parts = token.split(":");
    this.aws_access_key_id = parts[0];
    this.aws_secret_access_key = parts[1];

  }

  handleTokenChange = (token) => {
    if (!token) {
      return;
    }
    this.authCookie = token;
    Cookie.set('auth', token);
    this.parseAuthToken(token);
    this.forceUpdate();
  }
  waitForUpdated = (cb, slideToSave, times) => {
    times = times || 1;
    if (times > 10 ) {
      console.log("failed in time");
      return;
    }
    var handleResp = (err, data) => {
      var slideFromState = JSON.parse(data.Body);
      console.log("checking for consistent read", slideFromState.edited, slideToSave.edited)
      if (slideFromState.edited === slideToSave.edited) {
        cb(slideToSave);
      } else {
        setTimeout(() => {
          times += 1
          this.waitForUpdated(cb, slideToSave, times)
        }, 200);
      }
    };

    this.s3.getObject({
      Bucket: "art.rumproarious.com",
      Key: this.key,
    }, handleResp);
  }
  saveSlide = (e) => {
    e.preventDefault();
    var slideToSave = Object.assign({}, this.state.slide);
    if (!slideToSave.artist) {
      slideToSave.artist = null;
    }
    slideToSave.edited = (new Date()).toISOString();
    console.log("About to save", slideToSave)
    e.preventDefault();
    var params = {
      Body: JSON.stringify(slideToSave),
      Bucket: "art.rumproarious.com",
      Key: this.key,
      ACL: "public-read",
      ContentType: "application/json",
    };

    this.s3.putObject(params, (err, data) => {
      if (err) {
        console.log(err, err.stack);
        return;
      }

      this.waitForUpdated((data) => {
        this.setState((state) => {
          state.slide = data;
        })
      }, slideToSave);

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
            <Artist artist={this.state.slide.artist} onValueChange={bindThatUpdate(this, (state, artist) => state.artist = artist)} />
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
