import { UserData, SimpleResourceData } from '../utils/api-client';

export interface TestUsers {
  validUser: UserData;
  anotherValidUser: UserData;
  invalidUser: UserData;
}

export interface TestSimpleResources {
  validResource: SimpleResourceData;
  anotherValidResource: SimpleResourceData;
  updatedResource: SimpleResourceData;
  invalidResource: SimpleResourceData;
}

export interface ApiEndpoints {
  health: string;
  swagger: string;
  auth: {
    signup: string;
    login: string;
  };
  simple: {
    base: string;
    byId: (id: number | string) => string;
  };
}

export interface ExpectedResponses {
  health: {
    status: string;
  };
  authSuccess: {
    signup: {
      message: string;
    };
    login: {
      message: string;
    };
  };
  simpleSuccess: {
    create: {
      message: string;
    };
    getAll: {
      message: string;
    };
    getById: {
      message: string;
    };
    update: {
      message: string;
    };
    delete: {
      message: string;
    };
  };
  errors: {
    unauthorized: {
      status: number;
      error: string;
    };
    notFound: {
      status: number;
      error: string;
    };
    badRequest: {
      status: number;
    };
    conflict: {
      status: number;
    };
  };
}

export interface TestConfig {
  timeouts: {
    short: number;
    medium: number;
    long: number;
  };
  retries: {
    default: number;
    network: number;
  };
  delays: {
    short: number;
    medium: number;
    long: number;
  };
}

export const testUsers: TestUsers = {
  validUser: {
    email: 'valid.user@example.com',
    password: 'SecurePassword123!'
  },
  
  anotherValidUser: {
    email: 'another.user@example.com',
    password: 'AnotherSecure456!'
  },
  
  invalidUser: {
    email: 'invalid-email',
    password: '123'
  }
};

export const testSimpleResources: TestSimpleResources = {
  validResource: {
    name: 'Test Simple Resource'
  },
  
  anotherValidResource: {
    name: 'Another Test Resource'
  },
  
  updatedResource: {
    name: 'Updated Test Resource'
  },
  
  invalidResource: {
    name: '' // Invalid empty name
  }
};

export const apiEndpoints: ApiEndpoints = {
  health: '/health',
  swagger: '/swagger/index.html',
  auth: {
    signup: '/auth/signup',
    login: '/auth/login'
  },
  simple: {
    base: '/simple',
    byId: (id: number | string) => `/simple/${id}`
  }
};

export const expectedResponses: ExpectedResponses = {
  health: {
    status: 'OK'
  },
  
  authSuccess: {
    signup: {
      message: 'User created successfully'
    },
    login: {
      message: 'Login successful'
    }
  },
  
  simpleSuccess: {
    create: {
      message: 'Simple created successfully'
    },
    getAll: {
      message: 'Simples retrieved successfully'
    },
    getById: {
      message: 'Simple retrieved successfully'
    },
    update: {
      message: 'Simple updated successfully'
    },
    delete: {
      message: 'Simple deleted successfully'
    }
  },
  
  errors: {
    unauthorized: {
      status: 401,
      error: 'Unauthorized'
    },
    notFound: {
      status: 404,
      error: 'Not found'
    },
    badRequest: {
      status: 400
    },
    conflict: {
      status: 409
    }
  }
};

export const testConfig: TestConfig = {
  timeouts: {
    short: 5000,
    medium: 10000,
    long: 30000
  },
  
  retries: {
    default: 3,
    network: 5
  },
  
  delays: {
    short: 100,
    medium: 500,
    long: 1000
  }
};
