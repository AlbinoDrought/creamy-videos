// https://docs.cypress.io/api/introduction/api.html

describe('Full Video CRUD Flow', () => {
  it('Uploads a file', () => {
    cy.visit('/upload');
    cy.contains('.submit.button', 'Upload');
    cy.get('[name="title"]').invoke('val').should('be.empty');
    cy.get('[name="tags"]').invoke('val').should('eq', 'home');
    cy.get('[name="description"]').invoke('val').should('be.empty');

    // upload a file
    cy.fixture('doggo_waddling.mp4', 'base64').then((content) => {
      cy.get('[name="file"]').upload(content, 'doggo_waddling.mp4', 'video/mp4');
      cy.get('[name="title"]').invoke('val').should('eq', 'doggo_waddling.mp4');

      // submit upload
      cy.get('.submit.button').click();

      cy.url().should('contain', '/watch/');

      // assert ui
      cy.contains('[aria-label="Video Title"]', 'doggo_waddling.mp4');
      cy.contains('[aria-label="Video Tags"]', 'home');

      // check for buttons
      cy.contains('.download.button', 'Download');
      cy.contains('.edit.button', 'Edit');
      cy.contains('.delete.button', 'Delete');

      // edit video
      cy.get('.edit.button').click();

      cy.url().should('contain', '/edit');
      
      // assert ui
      cy.get('[name="title"]').invoke('val').should('eq', 'doggo_waddling.mp4');
      cy.get('[name="tags"]').invoke('val').should('eq', 'home');
      cy.get('[name="description"]').invoke('val').should('be.empty');
      cy.get('.submit.button').should('contain', 'Save');

      // change input fields, save
      cy.get('[name="title"]').clear().type('My Waddling Doggo');
      cy.get('[name="tags"]').type(',doggo');
      cy.get('[name="description"]').type('This is a short video of my doggo waddling.');
      cy.get('.submit.button').click();

      cy.url().should('contain', '/watch/');
      cy.contains('[aria-label="Video Title"]', 'My Waddling Doggo');
      cy.contains('[aria-label="Video Tags"]', 'home');
      cy.contains('[aria-label="Video Tags"]', 'doggo');
      cy.contains('[aria-label="Video Description"]', 'This is a short video of my doggo waddling.');

      // delete it
      cy.get('.delete.button').click();
      cy.contains('.delete.button', 'Click 3 more times to confirm');
      cy.get('.delete.button').click();
      cy.contains('.delete.button', 'Click 2 more times to confirm');
      cy.get('.delete.button').click();
      cy.contains('.delete.button', 'Click 1 more times to confirm');
      cy.get('.delete.button').click();

      // pseudo-check that we have returned to the home screen
      cy.title().should('contain', 'Home');
    });
  });
});
