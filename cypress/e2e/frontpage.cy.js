describe("Load front page", () => {
  it("Has title TripleWorks", () => {
    cy.visit("/");
    cy.get("h1").should("have.text", "TripleWorks");
  });

  it("can create new item", () => {
    cy.visit("/");
    cy.get("#type-select").select("ReportingGroup");
    cy.get("#new-btn").click();
    cy.get("#commit-input").type("create an empty reporting group");
    cy.intercept("POST", "/commit").as("commitReq");
    cy.get("#commit-btn").click();
    cy.wait("@commitReq").then((intersception) => {
      expect(intersception.response.statusCode).to.equal(200);
    });

    cy.get("#status-message").should("contain", "Successfully updated");
  });

  it("can download xiidm", () => {
    cy.request("/xiidm").then((resp) => {
      expect(resp.status).to.equal(200);
      expect(resp.body.length).to.be.greaterThan(10);
    });
  });

  it("can retrieve crpss border ptdfs", () => {
    cy.request("/cross-border-ptdfs").then((resp) => {
      expect(resp.status).to.equal(200);
    });
  });
});
