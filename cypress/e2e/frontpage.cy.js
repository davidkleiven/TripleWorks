describe("Load front page", () => {
  it("Has title TripleWorks", () => {
    cy.visit("/");
    cy.get("h1").should("have.text", "TripleWorks");
  });
});
