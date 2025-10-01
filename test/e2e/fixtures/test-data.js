// test-data.js - Test fixtures and sample data

const testUsers = {
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

const testSimpleResources = {
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

const apiEndpoints = {
  health: '/health',
  swagger: '/swagger/index.html',
  auth: {
    signup: '/auth/signup',
    login: '/auth/login'
  },
  simple: {
    base: '/simple',
    byId: (id) => `/simple/${id}`
  }
};

const expectedResponses = {
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

const testConfig = {
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

module.exports = {
  testUsers,
  testSimpleResources,
  apiEndpoints,
  expectedResponses,
  testConfig
};
