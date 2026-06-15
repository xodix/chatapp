export class Request {
  static async getMessages(
    request: GetMessagesRequest,
  ): Promise<GetMessagesResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    };

    const response = await fetch(
      `http://localhost:8080/chats/messages?id=${request.chat_id}`,
      requestOptions as RequestInit,
    );
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async createChat(request: CreateChatRequest): Promise<string> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "POST",
      headers: myHeaders,
      body: JSON.stringify(request),
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/chats/",
      requestOptions as RequestInit,
    );
    if (!response) {
      throw "Could not get the response";
    }

    return await response.text();
  }

  static async register(request: RegisterRequest): Promise<RegisterResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");

    const requestOptions = {
      method: "POST",
      headers: myHeaders,
      body: JSON.stringify(request),
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/user/register",
      requestOptions as RequestInit,
    ).catch((error) => {
      if (error.status === 401) {
        localStorage.removeItem("token");
      }
      console.error(error);
    });
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async login(request: LoginRequest): Promise<LoginResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");

    const requestOptions = {
      method: "POST",
      headers: myHeaders,
      body: JSON.stringify(request),
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/user/login",
      requestOptions as RequestInit,
    ).catch((error) => {
      if (error.status === 401) {
        localStorage.removeItem("token");
      }
      console.error(error);
    });
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async accountDetails(): Promise<DetailsResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/user/",
      requestOptions as RequestInit,
    ).catch((error) => {
      if (error.status === 401) {
        localStorage.removeItem("token");
      }
      console.error(error);
    });
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async updateUser(request: UpdateRequest): Promise<string> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "PUT",
      headers: myHeaders,
      body: JSON.stringify(request),
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/user/",
      requestOptions as RequestInit,
    );
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async getChats(): Promise<GetChatsResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/chats/",
      requestOptions as RequestInit,
    ).catch((error) => {
      if (error.status === 401) {
        localStorage.removeItem("token");
      }
      console.error(error);
    });
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }

  static async getAllUsers(): Promise<GetAllUsersResponse> {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");
    const token = localStorage.getItem("token");
    if (!token) {
      throw "No token found";
    }
    myHeaders.append("Authorization", `Bearer ${token}`);

    const requestOptions = {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    };

    const response = await fetch(
      "http://localhost:8080/user/all",
      requestOptions as RequestInit,
    ).catch((error) => {
      if (error.status === 401) {
        localStorage.removeItem("token");
      }
      console.error(error);
    });
    if (!response) {
      throw "Could not get the response";
    }

    return await response.json();
  }
}

export type RegisterRequest = {
  email: string;
  password: string;
  name: string;
  surname: string;
  birthdate: Date;
};

export type RegisterResponse = {
  token: string;
  user_id: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type LoginResponse = {
  token: string;
  user_id: string;
};

export type DetailsResponse = {
  id: string;
  email: string;
  name: string;
  surname: string;
  birthdate: Date;
};

export type UpdateRequest = {
  email: string;
  name: string;
  surname: string;
  birthdate: Date;
};

export type ChatMessage = {
  author: string;
  message: string;
  name: string;
  surname: string;
};

export type Chat = {
  id: string;
  name: string;
  messages: ChatMessage[];
  members: number[];
};

export type GetChatsResponse = {
  chats: Chat[];
};

export type GetAllUsersResponse = {
  users: User[];
};

export type User = {
  id: string;
  name: string;
  surname: string;
};

export type CreateChatRequest = {
  name: string;
  members: string[];
};

export type GetMessagesRequest = {
  chat_id: string;
};

export type GetMessagesResponse = ChatMessage[];
