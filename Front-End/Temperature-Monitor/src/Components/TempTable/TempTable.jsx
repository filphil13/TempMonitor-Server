import React, { Component } from 'react';
import {useState, useEffect} from 'react';

const API_URL = "https://temp-monitor-a38f32c02c5e.herokuapp.com"



var SensorList;

function TempTable() {
    
    const [SensorList, setSensorList] = useState([]);
    const [TableBodyHTML, setTableBodyHTML] = useState("");


    useEffect(() => 
    {
        setInterval(() => {
            getRecentData()
        }, 10000);
    })


    function getRecentData(){

        const recentSensorData = [];
        fetch(API_URL + "/api/recent")
        .then(async response => {
            const data = await response.json();

            // check for error response
            if (!response.ok) {
                // get error message from body` or default to response statusText
                const error = (data && data.message) || response.statusText;
                return Promise.reject(error);
            }
            else if(data==[]){
                return 
            }
            sensorList = recentSensorData;
        });

        var TABLEBODYHTML = <></>
        if (SensorList.length > 0){
            TABLEBODYHTML = SensorList.map((sensor) =>(
                <tr key={sensor.name} class="border-b border-gray-200 dark:border-gray-700">
                    <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800">
                        {sensor.Name}
                    </th>
                    <td class="px-6 py-4">
                        {String(sensor.Temperature)}
                    </td>
                    <td class="px-6 py-4 bg-gray-50 dark:bg-gray-800">
                        {String(sensor.Humidity)}
                    </td>
                    <td class="px-6 py-4">
                        {String(sensor.Time)} seconds ago
                    </td>
                </tr>
            ));
        }
        setTableBodyHTML(TableBodyHTML)
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