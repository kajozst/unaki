import {Events} from "@wailsio/runtime";






/////////////////////////////////////////////////////
// trigger command                     /////////////
///////////////////////////////////////////////////
const getElement = (id) => document.getElementById(id);
const command = () => {
    
    let eventData = ["nak", "req", "-k", "1", "-l", "5", "wss://e.nos.lol"];
  
    if (eventData) {
        Events.Emit({name: 'command', data: JSON.stringify(eventData)});
    } else {
        console.warn('cant send command');
    }
    
};

document.addEventListener('DOMContentLoaded', () => {

    if (getElement('command-button')) getElement('command-button').addEventListener('click', command);

});



///////////////////////////////////////////
// listener to a backend custom event
//const outElement = document.getElementById('output');


Events.On('commandRes', (event) => {
    const data = (event.data);
    processData(data)
});

function processData(data) {
    let nakEvents = [];
  
    try {
      // Split the string into individual object strings
      const objectStrings = data[0].split('{"kind":');
  
      // Loop through the object strings, parse each one, and add it to the nakEvents array
      for (let i = 1; i < objectStrings.length; i++) {
        const objectString = '{"kind":' + objectStrings[i];
        const obj = parseObject(objectString);
        nakEvents.push(obj);
      }
  
      //console.log("~~~~~~~");
      //console.log(nakEvents);
      //console.log("~~~~~~~");

      // NOW WE HAVE nakEvents AS AN ARRAY OF OBJECTS!!!WE SHOULD LOOP THROUGH!!! ((function renderEvents(nakEvents);))

      // Assuming you have a container element in your HTML where you want to render the objects
      const container = document.getElementById('container');

      // Loop through the nakEvents array and create HTML elements for each object
      nakEvents.forEach(event => {
        // Create a new div element to hold the event data
        const eventDiv = document.createElement('div');
        eventDiv.classList.add('event');
      
        // Create and append elements for each property of the event object
        const kindElement = document.createElement('p');
        kindElement.textContent = `Kind: ${event.kind}`;
        eventDiv.appendChild(kindElement);

        const contentElement = document.createElement('p');
        contentElement.textContent = `Content: ${event.content}`;
        eventDiv.appendChild(contentElement);
      
        const idElement = document.createElement('p');
        idElement.textContent = `ID: ${event.id}`;
        eventDiv.appendChild(idElement);
      
        const pubkeyElement = document.createElement('p');
        pubkeyElement.textContent = `Pubkey: ${event.pubkey}`;
        eventDiv.appendChild(pubkeyElement);
      
        const createdAtElement = document.createElement('p');
        createdAtElement.textContent = `Created At: ${event.created_at}`;
        eventDiv.appendChild(createdAtElement);
      
        // Render the tags as a list
        const tagsElement = document.createElement('ul');
        event.tags.forEach(tag => {
          const tagItem = document.createElement('li');
          tagItem.textContent = `[${tag[0]}] ${tag[1]}`;
          tagsElement.appendChild(tagItem);
        });
        eventDiv.appendChild(tagsElement);
      
        // Append the event div to the container
        container.appendChild(eventDiv);
      });




    } catch (error) {
      console.error("Error parsing JSON:", error);
    }
  }
  
  function parseObject(objectString) {
    try {
      return JSON.parse(objectString);
    } catch (error) {
      console.error("Error parsing object:", error);
      return null;
    }
  }
  








const timeElement = document.getElementById('time');

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});


