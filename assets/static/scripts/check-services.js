    const STATUS_DIV = document.getElementById('matrix-status');
    // The homeserver you want to check
   const HOMESERVER = 'https://tavern.borderlandkeep.com';

    // How often to re-check (in milliseconds)
    const CHECK_INTERVAL = 30000;

    async function checkMatrixStatus() {
     // STATUS_DIV.textContent = 'Checking...';
      STATUS_DIV.className = 'status-badge status-checking';

      try {
        // The /_matrix/client/versions endpoint is a lightweight,
        // unauthenticated endpoint that all Matrix homeservers expose
        const response = await fetch(`${HOMESERVER}/_matrix/client/versions`, {
          method: 'GET',
          // A short timeout avoids hanging for too long on a dead server
          signal: AbortSignal.timeout(5000),
        });

        if (response.ok) {
          const data = await response.json();
          const versions = data.versions.join(', ');
    //      STATUS_DIV.textContent = `Online — supported versions: ${versions}`;
          STATUS_DIV.className = 'status-badge status-online';
        } else {
  //        STATUS_DIV.textContent = `Degraded (HTTP ${response.status})`;
          STATUS_DIV.className = 'status-badge status-offline';
        }

      } catch (err) {
        // fetch() throws on network errors or if the timeout fires
        const reason = err.name === 'TimeoutError' ? 'timed out' : 'unreachable';
//        STATUS_DIV.textContent = `Offline (${reason})`;
        STATUS_DIV.className = 'status-badge status-offline';
      }
    }

    // Run once immediately, then on a schedule
    checkMatrixStatus();
    setInterval(checkMatrixStatus, CHECK_INTERVAL);
