import React from "react";
import { Card, Table, Container, Row, Col, Button } from "react-bootstrap";

function TableList() {
  let buttonStyles = {
    marginRight: "10px",
  };

  return (
    <>
      <Container fluid>
        <Row>
          <Col md="12">
            <Card className="strpied-tabled-with-hover">
              <Card.Header>
                <Card.Title as="h4">Flags</Card.Title>
                <Button variant="success" className="btn-fill float-right">
                  Create a new flag
                </Button>
              </Card.Header>
              <Card.Body className="table-full-width table-responsive px-0">
                <Table className="table-hover table-striped">
                  <thead>
                    <tr>
                      <th className="border-0">ID</th>
                      <th className="border-0">Flag</th>
                      <th className="border-0">Created At</th>
                      <th className="border-0">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>1</td>
                      <td>flag1.flag1.flag1</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                    <tr>
                      <td>2</td>
                      <td>flag2.flag2.flag2</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                    <tr>
                      <td>3</td>
                      <td>flag3.flag3.flag3</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                    <tr>
                      <td>4</td>
                      <td>flag4.flag4.flag4</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                    <tr>
                      <td>5</td>
                      <td>flag5.flag5.flag5</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                    <tr>
                      <td>6</td>
                      <td>flag6.flag6.flag6</td>
                      <td>2021-04-30T10:48:23+02:00</td>
                      <td>
                        <Button
                          variant="primary"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Update
                        </Button>
                        <Button
                          variant="danger"
                          className="btn-fill"
                          style={buttonStyles}
                        >
                          Delete
                        </Button>
                      </td>
                    </tr>
                  </tbody>
                </Table>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </>
  );
}

export default TableList;
