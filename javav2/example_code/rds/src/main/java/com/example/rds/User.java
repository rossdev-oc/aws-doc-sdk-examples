// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package com.example.rds;

public class User {

    private String username;
    private String  password;

    private String host;


    //getter
    public String getUsername(){
        return this.username;
    }

    public String getPassword(){
        return this.password;
    }

    String getHost(){
        return this.host;
    }

}
