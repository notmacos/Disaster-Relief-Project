const sqlite3 = require('sqlite3').verbose();
const express = require('express');
const app = express();
const fs = require('fs');
const dbPath = 'shellhacks.db';

// Add these lines near the top of the file, after the existing requires
const cors = require('cors');
app.use(cors());

// Check if the database file exists
if (!fs.existsSync(dbPath)) {
  console.error(`Database file not found: ${dbPath}`);
  console.log('Creating a new database file...');
  fs.writeFileSync(dbPath, '');
}

// Create a new database connection
const db = new sqlite3.Database(dbPath, sqlite3.OPEN_READWRITE, (err) => {
  if (err) {
    console.error('Error connecting to SQLite database:', err.message);
    console.log('Please ensure you have write permissions in this directory.');
    process.exit(1);
  }
  console.log('Connected to the SQLite database.');
  
  // Create the events table if it doesn't exist
  db.run(`CREATE TABLE IF NOT EXISTS events (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    eventType TEXT,
    location TEXT,
    time TEXT,
    poster TEXT,
    description TEXT,
    emergency BOOLEAN
  )`, (err) => {
    if (err) {
      console.error('Error creating events table:', err.message);
      return;
    }
    console.log('Events table ready.');
  });
});

// Modify the fetchEvents function to use a callback
function fetchEvents(callback) {
  db.all('SELECT * FROM events', [], (err, rows) => {
    if (err) {
      console.error('Error fetching events:', err.message);
      callback(err, null);
    } else {
      callback(null, rows);
    }
  });
}

// Function to add events to the 'events' table
function addEvent(eventType, location, time, poster, description, callback) {
  const query = 'INSERT INTO events (eventType, location, time, poster, description) VALUES (?, ?, ?, ?, ?)';
  db.run(query, [eventType, location, time, poster, description], function(err) {
    if (err) {
      console.error('Error adding event:', err.message);
      callback(err);
    } else {
      console.log('New event added to table.');
      callback(null);
    }
  });
}

// Function to remove events from the 'events' table
function removeEvent(id, callback) {
  db.run('DELETE FROM events WHERE ID = ?', [id], function(err) {
    if (err) {
      console.error('Error removing event:', err.message);
      callback(err);
    } else {
      console.log('Event removed from table.');
      callback(null);
    }
  });
}

// Create a new API endpoint to serve events data
app.get('/api/events', (req, res) => {
  fetchEvents((err, events) => {
    if (err) {
      res.status(500).json({ error: 'Error fetching events' });
    } else {
      res.json(events);
    }
  });
});

// API endpoint to add an event
app.post('/api/events', express.json(), (req, res) => {
  const { eventType, location, time, poster, description } = req.body;
  addEvent(eventType, location, time, poster, description, (err) => {
    if (err) {
      res.status(500).json({ error: 'Error adding event' });
    } else {
      res.json({ message: 'Event added successfully' });
    }
  });
});

// Add a new API endpoint for SOS events
app.post('/api/sos', express.json(), (req, res) => {
  const { eventType, location, description, name } = req.body;
  const time = new Date().toISOString();
  const emergency = true;

  addEvent(eventType, location, time, name, description, emergency, (err) => {
    if (err) {
      res.status(500).json({ error: 'Error adding SOS event' });
    } else {
      res.json({ message: 'SOS event added successfully' });
    }
  });
});

// API endpoint to remove an event
app.delete('/api/events/:id', (req, res) => {
  const id = req.params.id;
  removeEvent(id, (err) => {
    if (err) {
      res.status(500).json({ error: 'Error removing event' });
    } else {
      res.json({ message: 'Event removed successfully' });
    }
  });
});

// Start the server
const port = 1400;
app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
});

// Handle process termination
process.on('SIGINT', () => {
  db.close((err) => {
    if (err) {
      console.error('Error closing the database connection:', err.message);
    } else {
      console.log('Database connection closed.');
    }
    process.exit(0);
  });
});