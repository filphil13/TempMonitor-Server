import React from 'react';

export default function SensorCard({name, temperature, humidity, status, time}){
    
    
    return (
        <div
            className=' p-14 bg-red-400 flex justify-between shadow-md border border-gray-200 rounded-lg  dark:bg-gray-800 dark:border-gray-700"'
            onClick={() => {console.log("Clicked")}}>   

            <p>{name}</p>
            <p>{String(temperature)}</p>
            <p>{String(humidity)}</p>
            <p>{String(status)}</p>
            <p>{String(Math.floor(Date.now() / 1000) - time)} seconds ago</p>

        </div>
    );
};

