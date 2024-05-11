import React from 'react';
import { View, Text, TextInput, TouchableOpacity } from 'react-native';
import styles from '../styles/style';

const LoginForm = () => {

    const handleLogin = () => {
        // Your login logic here
        console.log('Login button pressed!');
    };


    return (
        <View style={styles.container}>
            <View style={styles.formContainer}>
                <Text style={styles.title}>daily track</Text>
                <TextInput
                    style={styles.input}
                    placeholder="email"
                    keyboardType="email-address"
                    autoCapitalize="none"
                    autoCorrect={false}
                />
                <TextInput
                    style={styles.input}
                    placeholder="password"
                    secureTextEntry
                    autoCapitalize="none"
                    autoCorrect={false}
                />
                <TouchableOpacity style={styles.button} onPress={handleLogin}>
                    <Text style={styles.buttonText}>Login</Text>
                </TouchableOpacity>
            </View>
        </View>
    );
};

export default LoginForm;