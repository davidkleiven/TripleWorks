describe("can apply json patch", () => {
  it("can apply json patch", () => {
    cy.visit("/patch-form");
    cy.get('input[type="file"]').selectFile(
      "cypress/fixtures/json_patch.json",
      { force: true },
    );
    cy.intercept("PATCH", "/resource").as("resource");
    cy.get('button[type="submit"]').click();

    cy.wait("@resource").then((inter) => {
      expect(inter.response.statusCode).to.equal(200);
    });
  });
});
