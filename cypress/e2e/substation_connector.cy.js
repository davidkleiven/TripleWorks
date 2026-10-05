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

  it("move sends correct body", () => {
    cy.visit("/");
    cy.get("#list-all-btn").click();
    cy.contains('button[hx-get*="ce8e57c7"]', "Con sub.").click();
    cy.get("h2").should("have.text", "Substation connector");

    cy.get("#from-substation-search-input").type("Substation A");
    cy.get("#to-substation-search-input").type("Substation B");
    cy.get("#from-substation-results").contains("span", "Substation A").click();
    cy.get("#to-substation-results").contains("span", "Substation B").click();
    cy.get("#from-substation-display").should("contain", "Substation A");
    cy.get("#to-substation-display").should("contain", "Substation B");

    cy.intercept("POST", "/move/**").as("moveReq");

    // Moving to the substations the line already is connected to does nothing
    cy.get("#move-substations-btn").click();
    cy.wait("@moveReq").then((intersception) => {
      let body = intersception.request.body;

      expect(body.split("substation-mrid=").length - 1).to.equal(2);
      expect(intersception.response.statusCode).to.equal(200);
    });
    cy.get("#status-message").should("contain", "Nothing to move");

    // Move both ends of the line to Substation C and Substation D
    cy.get("#from-substation-search-input").clear().type("Substation C");
    cy.get("#from-substation-results").contains("span", "Substation C").click();
    cy.get("#to-substation-search-input").clear().type("Substation D");
    cy.get("#to-substation-results").contains("span", "Substation D").click();

    cy.get("#from-substation-display").should("contain", "Substation C");
    cy.get("#to-substation-display").should("contain", "Substation D");

    cy.get("#move-substations-btn").click();
    cy.wait("@moveReq").then((intersception) => {
      expect(intersception.response.statusCode).to.equal(200);
    });
    cy.get("#status-message").should("contain", "Successfully moved");
  });
});
