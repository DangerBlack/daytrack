import React, { useState } from 'react';
import { View, Text, TextInput, Button, ActivityIndicator, Alert } from 'react-native';
import {Picker} from '@react-native-picker/picker';
import styles from '../styles/style';
import { load_key } from '../models/storage';
import { track_url } from '../models/endpoints';

async function add_track(set_is_loading, track_info) 
{
    console.log('Adding a track')
    set_is_loading(true);
    try {
        const api_key = await load_key();
        const response = await fetch(track_url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Key ${api_key}`
            },
            body: JSON.stringify(track_info),
        });
        console.log('Track added')
        const data = await response.json();
        console.log(data);

    } catch (error) {
      console.error(error);
    } finally {
      set_is_loading(false);
    }
}

const CreateNewTrack = () => 
{
    const [is_loading, set_is_loading] = useState(false);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [visibility, setVisibility] = useState('private');
    const [status, setStatus] = useState('enabled');

    const handleSubmit = () => {
        // Handle form submission
        Alert.alert('Form Submitted', `Name: ${name}\nDescription: ${description}\nVisibility: ${visibility}\nStatus: ${status}`);

        add_track(set_is_loading, {name, description, visibility, status}).catch(console.error);
    };

    if (is_loading) {
        return (
            <View style={styles.container}>
                <ActivityIndicator size="large" color="#007BFF" />
            </View>
        );
    }

    return (
        <View style={styles.tackContainer}>
            <Text style={styles.tackLabel}>Name</Text>
            <TextInput
                style={styles.trackInput}
                placeholder="Enter name"
                value={name}
                onChangeText={setName}
            />

            <Text style={styles.tackLabel}>Description</Text>
            <TextInput
                style={styles.trackInput}
                placeholder="Enter description"
                value={description}
                onChangeText={setDescription}
            />

            <Text style={styles.tackLabel}>Visibility</Text>
            <Picker
                selectedValue={visibility}
                style={styles.trackPicker}
                onValueChange={(itemValue) => setVisibility(itemValue)}
            >
                <Picker.Item label="Private" value="private" />
                <Picker.Item label="Public (Read Only)" value="public_r" />
                <Picker.Item label="Public (Read/Write)" value="public_rw" />
            </Picker>

            <Text style={styles.tackLabel}>Status</Text>
            <Picker
                selectedValue={status}
                style={styles.trackPicker}
                onValueChange={(itemValue) => setStatus(itemValue)}
            >
                <Picker.Item label="Enabled" value="enabled" />
                <Picker.Item label="Disabled" value="disabled" />
            </Picker>

            <Button title="Submit" onPress={handleSubmit} />
        </View>
    );
};


export default CreateNewTrack;