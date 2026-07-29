    const DARK = {
      label: 'dark',
      foreground:           '#a5a2a2',
      'foreground-bold':    '#a5a2a2',
      cursor:               '#a5a2a2',
      background:           '#090300',
      'color-black':        '#090300',
      'color-black-bright': '#5c5855',
      'color-red':          '#db2d20',
      'color-red-bright':   '#db2d20',
      'color-green':        '#01a252',
      'color-green-bright': '#01a252',
      'color-yellow':       '#fded02',
      'color-yellow-bright':'#fded02',
      'color-blue':         '#01a0e4',
      'color-blue-bright':  '#01a0e4',
      'color-magenta':      '#a16a94',
      'color-magenta-bright':'#a16a94',
      'color-cyan':         '#b5e4f4',
      'color-cyan-bright':  '#b5e4f4',
      'color-white':        '#a5a2a2',
      'color-white-bright': '#f7f7f7',
    };

    const LIGHT = {
      ...DARK,
      label: 'light',
      foreground:        '#4a4543',
      'foreground-bold': '#4a4543',
      cursor:            '#4a4543',
      background:        '#f7f7f7',
    };

    const SWATCH_KEYS = [
      { key: 'background',           name: 'Background' },
      { key: 'foreground',           name: 'Foreground' },
      { key: 'color-black',          name: 'Black' },
      { key: 'color-black-bright',   name: 'Black Bright' },
      { key: 'color-red',            name: 'Red' },
      { key: 'color-green',          name: 'Green' },
      { key: 'color-yellow',         name: 'Yellow' },
      { key: 'color-blue',           name: 'Blue' },
      { key: 'color-magenta',        name: 'Magenta' },
      { key: 'color-cyan',           name: 'Cyan' },
      { key: 'color-white',          name: 'White' },
      { key: 'color-white-bright',   name: 'White Bright' },
    ];

    const toggle     = document.getElementById('theme-toggle');
    const html       = document.documentElement;
    const hint       = document.getElementById('current-theme-hint');
    const labelDark  = document.getElementById('label-dark');
    const labelLight = document.getElementById('label-light');
    const palette    = document.getElementById('palette');
    const cssOutput  = document.getElementById('css-output');

    function renderPalette(theme) {
      palette.innerHTML = '';
      SWATCH_KEYS.forEach(({ key, name }) => {
        const hex = theme[key];
        palette.innerHTML += `
          <div class="swatch">
            <div class="swatch-color" style="background:${hex}"></div>
            <div class="swatch-info">
              <div class="swatch-name">${name}</div>
              <div class="swatch-hex">${hex}</div>
            </div>
          </div>`;
      });
    }

    function renderCSS(theme) {
      const groups = [
        { comment: 'Special', keys: ['foreground','foreground-bold','cursor','background'] },
        { comment: 'Black',   keys: ['color-black','color-black-bright'] },
        { comment: 'Red',     keys: ['color-red','color-red-bright'] },
        { comment: 'Green',   keys: ['color-green','color-green-bright'] },
        { comment: 'Yellow',  keys: ['color-yellow','color-yellow-bright'] },
        { comment: 'Blue',    keys: ['color-blue','color-blue-bright'] },
        { comment: 'Magenta', keys: ['color-magenta','color-magenta-bright'] },
        { comment: 'Cyan',    keys: ['color-cyan','color-cyan-bright'] },
        { comment: 'White',   keys: ['color-white','color-white-bright'] },
      ];

      let html = `<span class="token-special">:root</span> <span class="token-punc">{</span>\n`;
      groups.forEach(({ comment, keys }) => {
        html += `  <span class="token-comment">/* ${comment} */</span>\n`;
        keys.forEach(k => {
          const pad = ' '.repeat(Math.max(1, 26 - k.length));
          html += `  <span class="token-prop">--${k}</span><span class="token-punc">:${pad}</span><span class="token-value">${theme[k]}</span><span class="token-punc">;</span>\n`;
        });
      });
      html += `<span class="token-punc">}</span>`;
      cssOutput.innerHTML = html;
    }

    function applyTheme(isLight) {
      const theme = isLight ? LIGHT : DARK;
      html.setAttribute('data-theme', theme.label);
      // hint.textContent = theme.label + ' mode';
      labelDark.classList.toggle('active', !isLight);
      labelLight.classList.toggle('active', isLight);
//     renderPalette(theme);
//     renderCSS(theme);
    }

    toggle.addEventListener('change', () => applyTheme(toggle.checked));

    // Respect system preference on load
    const preferLight = window.matchMedia('(prefers-color-scheme: light)').matches;
    toggle.checked = preferLight;
    applyTheme(preferLight);
