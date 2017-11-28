// Navigation Functions

function toList() {
    $('#note-editor').hide();
    $('#note-viewer').hide();
    $('#note-list').show();
}

function toViewer() {
    $('#note-editor').hide();
    $('#note-list').hide();
    $('#note-viewer').show();
}

function toEditor() {
    $('#note-viewer').hide();
    $('#note-list').hide();
    $('#note-editor').show();
}

// Requests

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
      } else if (data.errors.length > 0) {
        // Display any errors.
        console.log(data.errors);
      }
    });
  };
}

function getNoteDeleteFunc(id) {
  return function(e) {
    // Confirm deletion.
    if (confirm("Are you sure you would like to delete the note? This action cannot be undone.")) {
      // Send delete request.
      $.ajax({
        url: getNoteUrl(id),
        type: 'DELETE',
        headers: {
          'Authorization': getAuthHeader()
        }
      }).done(function(data) {
        if (data.errors.length > 0) {
          // Log errors.
          console.log(errors);
        } else {
          // Update the table.
          updateNoteList();
        }
      });
    }
  };
}

function getNoteUpdateFunc(id) {
  return function(e) {
    e.preventDefault();

    // Hide help text under controls.
    $('#note-title-help').hide();
    $('#note-content-help').hide();

    // Send the data in a PUT request.
    $.ajax({
      url: getNoteUrl(id),
      type: 'PUT',
      data: {
        "title": $('#note-title-edit').val(),
        "content": $('#note-content-edit').val()
      },
      headers: {
        'Authorization': getAuthHeader()
      }
    }).done(function(data) {
      // Log errors.
      if (data.errors.length > 0) {
        console.log(data.errors);
      }

      // Show help text under controls where the data is invalid.
      if (data.fields["title"]) {
        $('#note-title-help').text(data.fields["title"]);
        $('#note-title-help').show();
      } else if (data.fields["content"]) {
        $('#note-content-help').text(data.fields["content"]);
        $('#note-title-help').show();
      } else {
        updateNoteList();
        toList();
      }
    });
  };
}

function getNoteEditFunc(id) {
  return function(e) {
    // Show the correct view.
    toEditor();

    // Change the submit button's value.
    $('#note-update').val('Update');

    // Query for the note's existing data.
    $.ajax({
      url: getNoteUrl(id),
      type: 'GET',
      headers: {
        'Authorization': getAuthHeader()
      }
    }).done(function(data) {
      // Fill out the form with existing data.
      if (data.models[0]) {
        var model = data.models[0];
        $('#note-title-edit').val(model.title);
        $('#note-content-edit').val(model.content.String);
      } else if (data.errors.length > 0) {
        console.log(data.errors);
      }

      // Set the correct function to store the data upon submit.
      $('#note-update').unbind('click');
      $('#note-update').on('click', getNoteUpdateFunc(id));
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
        var actions = $('<td></td>').appendTo(row);
        var showButton = $('<a href="#" class="btn">Show</a>').appendTo(actions);
        var editButton = $('<a href="#" class="btn">Edit</a>').appendTo(actions);
        var deleteButton = $('<a href="#" class="btn">Delete</a>').appendTo(actions);
        showButton.on('click', getNoteShowFunc(model.id));
        editButton.on('click', getNoteEditFunc(model.id));
        deleteButton.on('click', getNoteDeleteFunc(model.id));
      }
    } else if (data.errors) {
      for (var i = 0; i < data.errors.length; i++) {
        console.log(data.errors[i]);
      }
    }
  });
}

// Event handlers

$('#logout').on('click', function(e) {
  // Clear the access token.
  window.sessionStorage.accessToken = "";

  // Redirect to login page.
  window.location = '/login.html';
});

$('#tolist').on('click', function(e) {
  toList();
});

$('#newnote').on('click', function(e) {
  // Switch to the editor.
  toEditor();

  // Clear the editor values.
  $('#note-title-edit').val('');
  $('#note-content-edit').val('');

  // Change the value of the submit button.
  $('#note-update').val('Store');

  // Set the proper function for the submit button.
  $('#note-update').unbind('click');
  $('#note-update').bind('click', function(e) {
    e.preventDefault();

    // Hide help text under controls.
    $('#note-title-help').hide();
    $('#note-content-help').hide();

    // Send a request to store the new note.
    $.ajax({
      url: '/api/note',
      type: 'POST',
      data: {
        "title": $('#note-title-edit').val(),
        "content": $('#note-content-edit').val()
      },
      headers: {
        "Authorization": getAuthHeader()
      }
    }).done(function(data) {
      if (data.errors.length > 0) {
        // Log errors.
        console.log(data.errors);
      } else if (data.fields["title"]) {
        $('#note-title-help').text(data.fields["title"]);
        $('#note-title-help').show();
      } else if (data.fields["content"]) {
        $('#note-content-help').text(data.fields["content"]);
        $('#note-title-help').show();
      } else {
        // If successful, update the list and return to it.
        updateNoteList();
        toList();
      }
    });
  });
});

// Setup logic

updateNoteList();

$('#note-viewer').hide();
$('#note-editor').hide();
