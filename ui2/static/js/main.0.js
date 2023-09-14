/**
 * Submit the form containing the element
 * @param {HTMLElement} el
 */
window.cvSubmitNearestForm = function (el) {
  var form = null;
  var parent = el;
  for (var i = 0; i < 10; i += 1) {
    parent = parent.parentNode;
    if (!parent) {
      break;
    }
    if (parent.tagName === 'FORM') {
      form = parent;
      break;
    }
  }
  if (!form) {
    console.error('Element: ', el);
    throw new Error('No form found around element');
  }

  form.submit();
};

// when we're looking at a list of videos
// and we click one
// the video should automatically play
// -----
// the reality:
// 1. we click a video link
// 2. the user-interaction context is lost for the new page
// 3. the browser complains and doesn't play the video
// -----
// the solution: use SPA-style page loading.
// I tried using hx-boost (HTMX) originally,
// but I had some issues with videos continuing to play in the background,
// the page scrolling to the bottom weirdly, etc.
// the below code works great on my machine:
if (window.fetch) {
  window.addEventListener('popstate', function () {
    // when the user goes back, reload the page from the server
    window.location.reload();
  });
  document.querySelectorAll('a[cv-boost="true"]').forEach(function (el) {
    el.addEventListener('click', function (e) {
      e.preventDefault();
      const target = el.getAttribute('href');
      fetch(target)
        .then(function (resp) {
          return resp.text();
        })
        .then(function (text) {
          window.history.pushState({}, '', target);
          document.body.innerHTML = text;
        })
        .catch(function (ex) {
          // try to fallback to regular navigation if we break something
          console.error(ex);
          window.location.href = target;
          throw ex;
        });
    });
  });
}
