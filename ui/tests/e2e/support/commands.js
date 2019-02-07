// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add('upload', { prevSubject: 'element' }, (subject, file, fileName, mimeType) => {
  cy.window().then(window => {
    Cypress.Blob.base64StringToBlob(file, mimeType).then(blob => {
      const testFile = new window.File([blob], fileName, { type: mimeType });
      // hack to force set file
      Object.defineProperties(subject[0], {
        files: {
          configurable: true,
          value: [testFile],
        },
      });
      cy.wrap(subject).trigger('change');
    });
  });
});
