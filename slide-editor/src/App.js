import React, { Component } from 'react';
import { Navbar, Grid, Row, Col, Image, Button} from 'react-bootstrap';
import './App.css';
import Cookie from 'js-cookie';
import queryString from './query-string.js';
import {SiteTitle, SiteUrl} from './Site.js';
import {ArtistName, ArtsyURL, ArtistWikipediaURL} from './Artist.js';
import {Auth} from './Auth.js'
import S3 from 'aws-sdk/clients/s3';

// var SLIDE = {
//   site: {
//     title: "Imported",
//     url: "http://exampl.org"
//   },
//   page: {
//     title: "Dark freshness by Wassily KandinskySize: 19.7x26.1 cm",
//     url: "http://artist-kandinsky.tumblr.com/post/156438435180",
//     published: "2018-07-26T13:45:55.643125159-06:00",
//     GUIDHash: "054fd38619b0b3752d79ca3d66182796"
//   },
//   content: "\u003cimg src=\"http://68.media.tumblr.com/b9553f51978e311916d81066c7defad4/tumblr_okfouiuALn1vbru6ho1_500.jpg\"\u003e\u003cbr\u003e\u003cbr\u003e\u003cp\u003e\u003cstrong\u003e\u003ca href=\"https://goo.gl/TIsfkS\"\u003eDark freshness\u003c/a\u003e\u003c/strong\u003e by \u003ca href=\"http://artist-kandinsky.tumblr.com\"\u003eWassily Kandinsky\u003c/a\u003e\u003c/p\u003e\u003cbr\u003eSize: 19.7x26.1 cm",
//   guid_hash: "054fd38619b0b3752d79ca3d66182796.a98f339d491468f9372651f7f90e072ea3add478dea2075c8ba9181b038eef38",
//   source_image_url: "http://68.media.tumblr.com/b9553f51978e311916d81066c7defad4/tumblr_okfouiuALn1vbru6ho1_1280.jpg",
//   archived_image: {
//     url: "http://68.media.tumblr.com/b9553f51978e311916d81066c7defad4/tumblr_okfouiuALn1vbru6ho1_1280.jpg",
//     width: 496,
//     height: 650,
//     content_type: "image/jpg",
//     filename: "b26b253e494787267bdedd29c47024deebe240cd55bdadae43fcf6b6236d9df3.jpg"
//   }
// };

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
    console.log(parsed);
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

    var _this = this;
    this.key = parsed.data;
    window.fetch("http://art.rumproarious.com/" + parsed.data).then(function(resp) {
      return resp.json()
    }).then(function (data) {
      _this.setState((state) => {
        if (!data.artist) {
          data.artist = {};
        }
        state.slide = data;
        return state;
      })
    })
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
  saveSlide = (e) => {
    e.preventDefault();
    var slide = Object.assign({}, this.state.slide);
    if (!slide.artist) {
      slide.artist = null;
    }
    slide.edited = (new Date()).toISOString();
    console.log("About to save", JSON.stringify(slide))
    e.preventDefault();
    var params = {
      Body: JSON.stringify(slide),
      Bucket: "art.rumproarious.com",
      Key: this.key,
      ACL: "public-read",
      ContentType: "application/json",
    };

    this.s3.putObject(params, (err, data) => {
      if (err) {
        console.log(err, err.stack);
      } else {
        console.log(data);
      }
      this.setState((state) => {
        state.slide = slide;
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
                <h3>Artist</h3>
              </Col>
            </Row>
            <Row>
              <Col xs={12}>
                <ArtistName name={this.state.slide.artist.name} onValueChange={bindThatUpdate(this, (state, name) => state.slide.artist.name = name)}></ArtistName>
              </Col>
            </Row>
            <Row>
              <Col xs={12}>
                <ArtsyURL url={this.state.slide.artist.artsy_url} artist_name={this.state.slide.artist.name} onValueChange={bindThatUpdate(this, (state, url) => state.slide.artist.artsy_url = url)}></ArtsyURL>
              </Col>
            </Row>
            <Row>
              <Col xs={12}>
                <ArtistWikipediaURL url={this.state.slide.artist.wikipedia_url} artist_name={this.state.slide.artist.name} onValueChange={bindThatUpdate(this, (state, url) => state.slide.artist.wikipedia_url = url)}></ArtistWikipediaURL>
              </Col>
            </Row>

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
