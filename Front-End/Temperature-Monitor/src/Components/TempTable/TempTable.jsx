import React, { Component } from 'react';
import {useState, useEffect} from 'react';

const API_URL = "https://temp-monitor-a38f32c02c5e.herokuapp.com/sensor/recent"


var TestSensorList = [
    {"name":"test","temperature":0,"humidity":0,"time":0},
    {"name":"test","temperature":0,"humidity":0,"time":0},
    {"name":"test","temperature":0,"humidity":0,"time":0},
];

var TABLEBODYHTML;

function TempTable() {
    const [SensorList, setSensorList] = useState(TestSensorList);
    useEffect(()=> {
        getRecentData()
        CreateTempBlocks()
        setInterval(updateTempTable, 3000);
    })

    function updateTempTable(){
            getRecentData()
            CreateTempBlocks()
    }


    function CreateTempBlocks(data){
        TABLEBODYHTML = <></>
        var TempTable = SensorList
        if (SensorList.length > 0){
            TABLEBODYHTML = TempTable.map((sensor) =>(
                <tr key={sensor.name} class="border-b border-gray-200 dark:border-gray-700">
                    <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800">
                        {sensor.name}
                    </th>
                    <td class="px-6 py-4">
                        {String(sensor.temperature)}
                    </td>
                    <td class="px-6 py-4 bg-gray-50 dark:bg-gray-800">
                        {String(sensor.humidity)}
                    </td>
                    <td class="px-6 py-4">
                        {String(sensor.time)} seconds ago
                    </td>
                </tr>
            ));
        }
    
    }

    function getRecentData(){
        fetch(API_URL)
        .then(async response => {
            const data = await response.json();

            // check for error response
            if (!response.ok) {
                // get error message from body or default to response statusText
                const error = (data && data.message) || response.statusText;
                return Promise.reject(error);
            }

            this.setSensorList(data)
        })
        .catch(error => {
            console.error('There was an error!', error);
        });
            
    }



	return(	
        <div class="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <thead class="text-xs text-gray-700 uppercase dark:text-gray-400">
                    <tr>
                        <th scope="col" class="px-6 py-3 bg-gray-50 dark:bg-gray-800">
                            Name
                        </th>
                        <th scope="col" class="px-6 py-3">
                            Temperature
                        </th>
                        <th scope="col" class="px-6 py-3 bg-gray-50 dark:bg-gray-800">
                            Humidity
                        </th>
                        <th scope="col" class="px-6 py-3">
                            Time
                        </th>
                    </tr>
                </thead>
                <tbody>
                    {TABLEBODYHTML}
                </tbody>
            </table>
        </div>

    );
}
 
export default TempTable;