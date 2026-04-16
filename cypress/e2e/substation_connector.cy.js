describe("substation connector works", () => {
  it("connect sends correct body", () => {
    cy.visit("/");
    cy.get("#list-all-btn").click();
    cy.contains('button[hx-get*="ce8e57c7"]', "Con sub.").click();
    cy.get("h2").should("have.text", "Substation connector");

    cy.get("#from-substation-search-input").type("Substation A");
    cy.get("#to-substation-search-input").type("Substation B");

    cy.get("#from-substation-results").should("not.be.empty");
    cy.get("#to-substation-results").should("not.be.empty");

    cy.get("#from-substation-results").find("span").first().click();
    cy.get("#to-substation-results").find("span").first().click();

    cy.get("#from-substation-display").should("not.contain", "No selection");
    cy.get("#to-substation-display").should("not.contain", "No selection");

    cy.intercept("POST", "/connect/**").as("connectionReq");

    cy.get("#connect-substations-btn").click();
    cy.wait("@connectionReq").then((intersception) => {
      let body = intersception.request.body;

      expect(body).to.contain("modelId=");
      expect(body.split("substation-mrid=").length - 1).to.equal(2);
      expect(intersception.response.statusCode).to.equal(200);
    });

    cy.get("#status-message").should("contain", "Successfully committed");
  });
});
