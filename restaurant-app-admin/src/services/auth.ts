export const getToken = () => {
    return localStorage.getItem("token");
};

export const isLoggedIn = () => {
    return !!getToken();
};

export const logout = () => {
    localStorage.removeItem("token");
};

export const verifyToken = async () => {
    const token = localStorage.getItem("token");

    if (!token) return false;

    try {
        const response = await fetch(
            "http://localhost:8080/api/auth/me",
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );

        return response.ok;
    } catch {
        return false;
    }
};