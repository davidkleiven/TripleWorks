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

  it("triggers connect-dangling lines on click", () => {
    cy.visit("/patch-form");
    cy.intercept("POST", "/connect-dangling").as("resource");
    cy.get("#connect-dangling-lines-btn").click();
    cy.wait("@resource").then((inter) => {
      expect(inter.response.statusCode).to.equal(200);
      expect(inter.response.headers["content-type"]).to.equal("text/html");
      expect(inter.response.body).to.contain("Inserted");
    });
  });
});
