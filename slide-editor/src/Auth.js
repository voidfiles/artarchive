import React, { Component } from 'react';
import { Navbar, Grid, Row, Col, FormGroup, ControlLabel, FormControl, Button} from 'react-bootstrap';

class Auth extends Component {
  handleTokenChange = (e) => {
    this.nextAuthToken = e.target.value;
    console.log("handleTokenChange", this.nextAuthToken);
  }

  handleSubmit = (e) => {
    this.props.onToken(this.nextAuthToken);
    e.preventDefault();
    console.log("handleSubmit");
  }

  render() {

    return (
      <div className="App">
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="#home">Slide Editor</a>
            </Navbar.Brand>
          </Navbar.Header>
        </Navbar>
        <Grid>
          <Row>
            <Col xs={12}>
              <form onSubmit={this.handleSubmit}>
                <FormGroup>
                  <ControlLabel>Auth Token</ControlLabel>
                  <FormControl
                    type="text"
                    placeholder="Enter text"
                    onChange={this.handleTokenChange}
                  />
                </FormGroup>
                <Button type="submit">Submit</Button>
              </form>
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

export {Auth};
