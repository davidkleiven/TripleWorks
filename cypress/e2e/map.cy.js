describe("test map", () => {
  it("can add and remove production", () => {
    cy.visit("/map");
    cy.window().then((win) => win.substationMarkers[0].openPopup());
    cy.intercept("POST", "/flow").as("flowReq");
    cy.contains("button", "Produce").click();

    // Only counts active
    cy.get('#active-production-form input[type="number"]').should(
      "have.length",
      1,
    );
    cy.get('#active-production-form input[type="hidden"]').should(
      "have.length",
      1,
    );
    cy.wait("@flowReq").then((inter) => {
      expect(inter.response.statusCode).to.equal(200);
      const body = inter.response.body;
      const lineMrid = "4e832836-ef53-458e-9711-903982551fcf";
      expect(body.flow).to.have.property(lineMrid);
    });

    // Check that that a flow label exists
    cy.get(".flow-value").should("exist");
    cy.get("#active-production-form").find("button").click();
    cy.get(".flow-value").should("not.exist");
    cy.get('#active-production-form input[type="number"]').should(
      "have.length",
      0,
    );
    cy.get('#active-production-form input[type="hidden"]').should(
      "have.length",
      0,
    );
  });
});
