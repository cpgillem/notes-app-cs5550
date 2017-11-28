// Core data functions

// Concatenates an authorization header for any requests to use.
function getAuthHeader() {
  return "Bearer " + window.sessionStorage.accessToken;
}

// Generates a URL to serve as the endpoint for a note.
function getNoteUrl(id) {
  return "/api/note/" + id;
}

// Generates a handler function for clicking on the Show button of a note.
function getNoteShowFunc(id) {
  return function(e) {
    // Show the correct view.
    $('#note-list').hide();
    $('#note-editor').hide();
    $('#note-viewer').show();
    
    // Request the note's data.
    $.ajax({
      url: getNoteUrl(id),
      type: 'GET',
      headers: {
        'Authorization': getAuthHeader()
      }
    }).done(function(data) {
      if (data.models[0]) {
        console.log(data.models[0]);
        // Display the note's data.
        $('#note-title').text(data.models[0].title);
        $('#note-date').text(data.models[0].date);
        $('#note-content').text(data.models[0].content.String);
      } else if (data.errors) {
        // Display any errors.
        console.log(data.errors);
      }
    });
  };
}

// Updates the main list of notes.
function updateNoteList() {
  $.ajax({
    url: '/api/note',
    type: 'GET',
    headers: {
      "Authorization": getAuthHeader()
    }
  }).done(function(data) {
    // Clear the table.
    $('#note-table-body').empty();

    // Load the table with notes.
    if (data.models) {
      for (var i = 0; i < data.models.length; i++) {
        console.log(data.models[i]);
        var model = data.models[i];
        
        // Add a table row.
        var row = $('<tr></tr>').appendTo('#note-table-body');
        row.append('<td>' + model.id + '</td>');
        row.append('<td>' + model.title + '</td>');
        row.append('<td>' + model.date + '</td>');
        var actions = $('<td></td>').appendTo(row);
        var editButton = $('<a href="#" class="btn">Show</a>').appendTo(actions);
        editButton.on('click', getNoteShowFunc(model.id));
      }
    } else if (data.errors) {
      for (var i = 0; i < data.errors.length; i++) {
        console.log(data.errors[i]);
      }
    }
  });
}

function createNote(id) {

}

function storeNote(id) {

}

function showNote(id) {

}

function editNote(id) {

}

function updateNote(id) {

}

function deleteNote(id) {

}

// Event handlers

$('#logout').on('click', function(e) {
  // Clear the access token.
  window.sessionStorage.accessToken = "";

  // Redirect to login page.
  window.location = '/login.html';
});

$('#tolist').on('click', function(e) {
  // Return to viewing the notes list.
  $('#note-list').show();
  $('#note-viewer').hide();
  $('#note-editor').hide();
});

// Setup logic

updateNoteList();

$('#note-viewer').hide();
$('#note-editor').hide();
