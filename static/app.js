const Controller = {
  // The search method is triggered when the form is submitted.
  // It prevents the default form submission, gets the search query from the form data,
  // and sends a fetch request to the server to search for the query.
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));

    const response = fetch(`/search?q=${data.query}`, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    }).then((response) => {
      // When the response is received, it is converted to JSON format
      // and passed to the updateCards method for rendering.
      response.json().then((results) => {
        if (results.length === 0) {
          const content = `
            <p>Sorry, we couldn't find any results.</p>
          `;
          createAndOpenModal(content);
        } else {
          Controller.updateCards(results);
        }
      });
    });
  },

  // The updateCards method takes an array of search results and generates
  // HTML code for each result to be displayed as a card on the webpage.
  // It then updates the card container with the generated HTML.
  updateCards: (results) => {
    const container = document.getElementById("card-container");
    const cards = [];
    for (let result of results) {
      const card = `
        <div class="col s12 m6 l4">
          <div class="card blue-grey darken-1">
            <div class="card-content white-text">
              <span class="card-title">${result.Title}</span>
              <p><strong>Character:</strong> ${result.Player}</p>
              <p><strong>Quote:</strong> ${result.Quote}</p>
              <p><strong>Act/Scene/Line:</strong> ${result.ActSceneLine}</p>
            </div>
            <div class="card-action">
              <a class="more-button" data-quote="${result.Quote}">More</a>
            </div>
          </div>
        </div>
      `;
      cards.push(card);
    }
    container.innerHTML = cards.join("");

    // Add event listeners to each "More" button so that when clicked,
    // it triggers the showContext method to display additional text surrounding the quote.
    const moreButtons = document.getElementsByClassName("more-button");
    for (let button of moreButtons) {
      button.addEventListener("click", Controller.showContext);
    }
  },
  
  // This function is called when a user clicks on the "More" button in a search result card.
  showContext: (ev) => {
    // Get the quote from the "data-quote" attribute of the button that was clicked.
    const quote = ev.target.dataset.quote;

    // Send a request to the server to get the context for the quote.
    const contextResponse = fetch(`/search-context?q=${quote}`, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    }).then((response) => {

      // Once the server responds with the context, display it in a modal.
      response.json().then((result) => {
        if (result.length === 0) {
          const content = `
            <p>Sorry, we couldn't find any results.</p>
          `;
          createAndOpenModal(content);
        } else {
          // If context was found, format it with the search term highlighted and display it in a modal.
          const context = result[0].Context;
          const searchTerm = document.querySelector("#query").value;
          const boldedSearchTerm = `<highlight>${searchTerm}</highlight>`;
          const boldedContext = context.replaceAll(searchTerm, boldedSearchTerm);
          const content = `
            <p>...${boldedContext}...</p>
          `;
          createAndOpenModal(content);
        }
      });
    });
  },
};  

// This function creates a modal with the given content and opens it.
const createAndOpenModal = (content) => {
  const modal = `
    <div class="modal">
      <div class="modal-content">
        ${content}
      </div>
      <div class="modal-footer">
        <a href="#!" class="modal-close waves-effect waves-green btn-flat">Close</a>
      </div>
    </div>
  `;

  // Checks if a modal is already open and closes it if it is.
  const existingModal = document.querySelector('.modal');
  if (existingModal) {
    existingModal.remove();
  }

  // Add the modal to the DOM and open it.
  document.body.insertAdjacentHTML('beforeend', modal);
  const modalElement = document.querySelector('.modal');
  const modalInstance = M.Modal.init(modalElement);
  modalInstance.open();
};

// Get the search form element and add a submit event listener to it.
const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
