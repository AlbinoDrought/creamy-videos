// https://docs.cypress.io/api/introduction/api.html

const uploadVideo = (title, tags, description = 'not an empty string', originalFileName = 'doggo_waddling.mp4') => cy.wrap(new Promise((resolve) => {
  cy.visit('/upload');
  cy.contains('.submit.button', 'Upload');
  cy.get('body [name="title"]').invoke('val').should('be.empty');
  cy.get('body [name="tags"]').invoke('val').should('eq', 'home');
  cy.get('body [name="description"]').invoke('val').should('be.empty');

  cy.get('body [name="file"]').attachFile({ filePath: 'doggo_waddling.mp4', encoding: 'binary' });

  cy.get('body [name="title"]').clear().type(title);
  cy.get('body [name="tags"]').clear().type(tags);
  cy.get('body [name="description"]').clear().type(description);
  cy.get('.submit.button').click();

  cy.url().should('contain', '/watch/');
  cy.get('video')
    .should('have.prop', 'duration')
    .and('be.greaterThan', '0');

  resolve();
}));

const assertVideoIsSeen = (title) => {
  cy.contains('[data-e2e="Video Thumbnail"]', title);
};

const searchForText = (term) => {
  cy.get('[data-e2e="Search"]').clear().type(term).type('{enter}');
  cy.url().should('contain', term);
};

const searchForTag = (tag) => {
  // no current tag search option on UI
  cy.visit(`/search?mode=tags&tags=${tag}`);
};

describe('Video Search Flow', () => {
  it('Searches videos', () => {
    const videos = [
      {
        title: 'foo',
        tags: 'home,foo',
      },
      {
        title: 'bar',
        tags: 'bar',
      },
      {
        title: 'foo bar',
        tags: 'foo,bar',
      },
      {
        title: 'baz',
        tags: 'home,baz,outsider',
      },
    ];

    // create all the videos
    let uploadPromise = cy.wrap(Promise.resolve());
    videos.forEach((video) => {
      uploadPromise = uploadPromise.then(() => uploadVideo(video.title, video.tags));
    });

    uploadPromise.then(() => {
      cy.visit('/');

      // we should see foo and baz on the home page
      assertVideoIsSeen('foo');
      assertVideoIsSeen('baz');

      searchForText('foo');
      assertVideoIsSeen('foo');
      assertVideoIsSeen('foo bar');

      searchForText('bar');
      assertVideoIsSeen('bar');
      assertVideoIsSeen('foo bar');

      searchForText('ba');
      assertVideoIsSeen('baz');
      assertVideoIsSeen('bar');
      assertVideoIsSeen('foo bar');

      searchForTag('foo');
      assertVideoIsSeen('foo');

      searchForTag('baz');
      assertVideoIsSeen('baz');

      searchForTag('bar');
      assertVideoIsSeen('bar');
      assertVideoIsSeen('foo bar');

      searchForTag('outsider');
      assertVideoIsSeen('baz');
    });
  });
});
