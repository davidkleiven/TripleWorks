describe("simple line page", () => {
  it("creates a line and reports success", () => {
    cy.visit("/simple");
    cy.get("#status-bar").should("have.text", "Idle");

    cy.intercept("POST", "/upload/lines*").as("uploadLine");

    cy.get("#from-input").type("Substation A");
    cy.get("#to-input").type("Substation B");
    cy.get("#length-input").type("10");
    cy.get("#voltage-input").type("400");
    cy.get("#submit-line").click();

    cy.wait("@uploadLine").then((interception) => {
      expect(interception.response.statusCode).to.equal(200);
      expect(interception.request.body).to.contain('"length":10');
      expect(interception.request.body).to.contain('"voltage":400');
    });

    cy.get("#status-bar").should("contain", "Successfully");
  });
});
