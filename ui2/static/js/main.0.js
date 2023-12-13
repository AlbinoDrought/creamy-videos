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

/**
 * Replace the content of a page entirely
 * @param {string} html 
 */
window.cvReplacePage = function (html) {
  var doc = document.open();
  doc.write(html);
  doc.close();
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
window.addEventListener('popstate', function (e) {
  // when the user goes back, reload the page from the server
  if (e.state && e.state.scrollY) {
    window.sessionStorage.setItem('cvBoostScrollY', e.state.scrollY);
  }
  window.location.reload();
});
var cvBoostScrollY = window.sessionStorage.getItem('cvBoostScrollY');
if (cvBoostScrollY) {
  cvBoostScrollY = parseInt(cvBoostScrollY, 10);
  window.sessionStorage.removeItem('cvBoostScrollY');
  window.scroll({ top: cvBoostScrollY });
}
window.cvPerformBind = function () {
  // fix above user-interaction issue
  document.querySelectorAll('a[cv-boost="true"]').forEach(function (el) {
    if (el.cvBoundBoost) {
      return;
    }
    el.cvBoundBoost = true;

    el.addEventListener('click', function (e) {
      e.preventDefault();
      var target = el.getAttribute('href');
      fetch(target)
        .then(function (resp) {
          return resp.text();
        })
        .then(function (text) {
          window.history.replaceState({ scrollY: window.scrollY }, '');
          window.history.pushState({}, '', target);
          window.cvReplacePage(text);
        })
        .catch(function (ex) {
          // try to fallback to regular navigation if we break something
          console.error(ex);
          window.location.href = target;
          throw ex;
        });
    });
  });

  // click a button multiple times to submit a form
  document.querySelectorAll('a[cv-confirm]').forEach(function (el) {
    if (el.cvBoundConfirm) {
      return;
    }
    el.cvBoundConfirm = true;

    var clicks = 0;
    var requiredClicks = 4;
    var timeoutHandle;
    var originalText = el.innerText;

    var updateText = function () {
      if (clicks === 0) {
        el.innerText = originalText;
      } else {
        var remainingClicks = requiredClicks - clicks;
        el.innerText = [
          'Click',
          remainingClicks,
          'more',
          remainingClicks === 1 ? 'time' : 'times',
          'to confirm',
        ].join(' ');
      }
    };

    var resetClicks = function () {
      clicks = 0;
      timeoutHandle = false;
      updateText();
    };

    var target = document.querySelector(el.getAttribute('cv-confirm'));
    if (!target) {
      console.error('cv-confirm element target not found!', el);
      return;
    }
    el.addEventListener('click', function (e) {
      e.preventDefault();
      clearTimeout(timeoutHandle);
      clicks += 1;
      if (clicks >= requiredClicks) {
        clicks = 0;
        updateText();
        target.submit();
      } else {
        updateText();
        timeoutHandle = setTimeout(resetClicks, 2000);
      }
    });
  });

  // if a file input is filled, default a text input to the file's name
  document.querySelectorAll('input[type="file"][cv-filename-default-to]').forEach(function (el) {
    if (el.cvBoundFilenameDefault) {
      return;
    }
    el.cvBoundFilenameDefault = true;

    var target = document.querySelector(el.getAttribute('cv-filename-default-to'));
    if (!target) {
      console.error('cv-filename-default-to element target not found!', el);
      return;
    }
    el.addEventListener('change', function (e) {
      var file = e.target.files[0];
      if (!target.value) {
        target.value = file.name;
      }
    });
  });

  document.querySelectorAll('[cv-infinite-scroll]').forEach(function (el) {
    if (el.cvBoundInfiniteScroll) {
      return;
    }
    el.cvBoundInfiniteScroll = true;

    var checkAlmostScrolledIntoView = function () {
      if (window.scrollY < el.getBoundingClientRect().top - 250) {
        return;
      }
      document.removeEventListener('scroll', checkAlmostScrolledIntoView);
      
      var target = el.getAttribute('cv-infinite-scroll');
      if (!target) {
        return;
      }
      fetch(target)
        .then(function (resp) {
          return resp.text();
        })
        .then(function (text) {
          var domParser = new DOMParser();
          var targetDocument = domParser.parseFromString(text, 'text/html');
          var targetContainer = targetDocument.querySelector('[cv-infinite-scroll-data]');
          if (!targetContainer) {
            throw new Error('Fetched target ' + target + ' but found no [cv-infinite-scroll-data]');
          }
          el.after(...targetContainer.children);
          // re-bind newly inserted elements
          window.cvPerformBind();
        })
    };
    document.addEventListener('scroll', checkAlmostScrolledIntoView);
  });
};
window.cvPerformBind();

console.log(`
# AlbinoDrought/creamy-videos

Repo: https://github.com/AlbinoDrought/creamy-videos

Source: https://${window.location.host}/source.tar.gz
`);
