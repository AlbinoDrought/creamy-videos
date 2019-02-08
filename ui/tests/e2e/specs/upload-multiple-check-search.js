// https://docs.cypress.io/api/introduction/api.html

const uploadVideo = (title, tags, description = 'not an empty string', originalFileName = 'doggo_waddling.mp4') => new Promise((resolve) => {
  cy.visit('/upload');
  cy.contains('.submit.button', 'Upload');
  cy.get('[name="title"]').invoke('val').should('be.empty');
  cy.get('[name="tags"]').invoke('val').should('eq', 'home');
  cy.get('[name="description"]').invoke('val').should('be.empty');

  cy.fixture('doggo_waddling.mp4', 'base64').then((content) => {
    cy.get('[name="file"]').upload(content, originalFileName, 'video/mp4');

    cy.get('[name="title"]').clear().type(title);
    cy.get('[name="tags"]').clear().type(tags);
    cy.get('[name="description"]').clear().type(description);
    cy.get('.submit.button').click();

    cy.url().should('contain', '/watch/');
    resolve();
  });
});

const assertVideoIsSeen = (title) => {
  cy.contains('[aria-label="Video Thumbnail"]', title);
};

const searchForText = (term) => {
  cy.get('[aria-label="Search"]').clear().type(term).type('{enter}');
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
    let uploadPromise = Promise.resolve();
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
